package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql" //blank import
	"gopkg.in/ldap.v3"
)

var jwtKey = []byte("my_secret_key")

// local struct for implementing Client interface
type client struct {
	model.Config
}

var usr string

// Auth implementation for the Client interface
func (c client) Auth(username, password string) error {
	// establish connection
	conn, err := connect(c.Host)
	if err != nil {
		WriteLogFile(err)
		return err
	}
	defer conn.Close()

	// perform initial read-only bind
	if err = conn.Bind(c.ROUser.Name, c.ROUser.Password); err != nil {
		WriteLogFile(err)
		return err
	}

	// find the user attempting to login
	results, err := conn.Search(ldap.NewSearchRequest(
		c.BaseDN, ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false, fmt.Sprintf("(%v=%v)", c.Filter, username),
		[]string{}, nil,
	))
	if err != nil {
		WriteLogFile(err)
		return err
	}
	if len(results.Entries) < 1 {
		return errors.New("not found")
	}

	// attempt auth
	log.Println(results.Entries)
	//role = strings.Join(results.Entries[0].Attributes[3].Values, "")
	return conn.Bind(results.Entries[0].DN, password)
}

// New client with the provided config
// If the configuration provided is invalid,
// or client is unable to connect with the config
// provided, an error will be returned
func New(config model.Config) (model.Client, error) {
	config, err := validateConfig(config)
	if err != nil {
		WriteLogFile(err)
		return nil, err
	}
	c := client{config}
	conn, err := connect(c.Host) // test connection
	if err != nil {
		WriteLogFile(err)
		return nil, err
	}
	if err = conn.Bind(c.ROUser.Name, c.ROUser.Password); err != nil {
		WriteLogFile(err)
		return nil, err
	}
	conn.Close()
	return c, err
}

// Helper functions

// establishes a connection with an ldap host
// (the caller is expected to Close the connection when finished)
func connect(host string) (*ldap.Conn, error) {
	c, err := net.DialTimeout("tcp", host, time.Second*8)
	if err != nil {
		WriteLogFile(err)
		return nil, err
	}
	conn := ldap.NewConn(c, false)
	conn.Start()
	return conn, nil
}

func validateConfig(config model.Config) (model.Config, error) {
	if config.BaseDN == "" || config.Host == "" || config.ROUser.Name == "" || config.ROUser.Password == "" {
		return model.Config{}, errors.New("[CONFIG] The config provided could not be validated")
	}
	if config.Filter == "" {
		config.Filter = "sAMAccountName"
	}
	return config, nil
}

func token(w http.ResponseWriter, username string) {
	var error model.Error
	db := database.DbConn()
	defer db.Close()
	if Role == "project manager" {
		selDB, err := db.Query("SELECT project_manager_email from project_manager WHERE project_manager_email =? ", username)
		defer selDB.Close()
		if err != nil {
			WriteLogFile(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if selDB.Next() == false {
			w.WriteHeader(http.StatusUnauthorized)
			error.Message = "Unauthorized user role"
			json.NewEncoder(w).Encode(error)
			return
		}
	} else if Role == "program manager" {
		selDB, err := db.Query("SELECT program_manager_email from program_manager WHERE program_manager_email =? ", username)
		defer selDB.Close()
		if err != nil {
			WriteLogFile(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if selDB.Next() == false {
			w.WriteHeader(http.StatusUnauthorized)
			error.Message = "Unauthorized user role"
			json.NewEncoder(w).Encode(error)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		error.Message = "Unauthorized user role"
		json.NewEncoder(w).Encode(error)
		return
	}
	expirationTime := time.Now().Add(3600 * time.Second).Unix()
	claims := &model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		WriteLogFile(err)
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ_-&" +
		"abcdefghijklmnopqrstuvwxyz%*" +
		"0123456789")
	length := 20
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	var jwtToken model.JwtToken
	jwtToken.AccessToken = tokenString
	jwtToken.TokenType = "bearer"
	jwtToken.Expiry = "3600"
	jwtToken.RefreshToken = b.String()
	json.NewEncoder(w).Encode(jwtToken)
	w.WriteHeader(http.StatusCreated)
	createdAt := time.Now()
	var query string = "Insert into token(username,access_token,expiration,role,created_at) values (?,?,?,?,?)"
	insert, err := db.Prepare(query)
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	insert.Exec(username, tokenString, expirationTime, Role, createdAt.Format("2006-01-02"))
	defer insert.Close()

	query = "Insert into refresh_token(username,access_token,refresh_token,created_at) values (?,?,?,?)"
	insert1, err := db.Prepare(query)
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	insert1.Exec(username, tokenString, b.String(), createdAt.Format("2006-01-02"))
	defer insert1.Close()
}

// SignIn : for user sign-in through LDAP
func (c *Commander) SignIn(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	var err error
	var error model.Error
	db := database.DbConn()
	defer db.Close()
	// create a new client
	if client, err = New(model.Config{
		BaseDN: "DC=sls,DC=ads,DC=valuelabs,DC=net",
		//BaseDN: "cn=ltest,ou=SERVICE ACCOUNTS,ou=SLS,dc=SLS,dc=ads,dc=valuelabs,dc=net",
		Filter: "userPrincipalName",
		ROUser: model.User{Name: "L test", Password: "Welcome@123"},
		Title:  "title",
		Host:   "10.10.52.113:389",
	}); err != nil {
		WriteLogFile(err)
		fmt.Println(err)
		return
	}
	var creds model.Credentials
	//	var pass string
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		WriteLogFile(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// var usr = creds.Username
	// var bytePassword = []byte(creds.Password)
	username := creds.Username
	password := creds.Password
	Role = creds.Role
	splitUser := strings.Split(username, "@")
	print := splitUser[0]
	user1 := fmt.Sprintf("%s@valuelabs.com", print)
	user2 := fmt.Sprintf("%s@sls.ads.valuelabs.net", print)
	err = client.Auth(user2, password)
	if err == nil {
		fmt.Println("Success!")
		token(w, user1)
	} else if err.Error() == "not found" {
		fmt.Println("H2")
		if errr := client.Auth(user1, password); errr != nil {
			fmt.Println("H3")
			WriteLogFile(errr)
			w.WriteHeader(http.StatusUnauthorized)
			error.Code = "401"
			error.Message = "Invalid Username or Password"
			json.NewEncoder(w).Encode(error)
			return
		} //else {
		fmt.Println("Success!")
		token(w, user1)
		//}
	} else {
		fmt.Println("H4")
		WriteLogFile(err)
		w.WriteHeader(http.StatusUnauthorized)
		error.Code = "401"
		error.Message = "Invalid Username or Password"
		json.NewEncoder(w).Encode(error)
		return
	}
}

// Refresh : to generate refresh tokens
func (c *Commander) Refresh(w http.ResponseWriter, r *http.Request) {
	fmt.Println("123")
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	claims := &model.Claims{}
	type refresh struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	fmt.Println("123")
	var ref refresh
	err := json.NewDecoder(r.Body).Decode(&ref)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		WriteLogFile(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var refToken = ref.RefreshToken
	var accToken = ref.AccessToken
	fmt.Println(refToken)
	fmt.Println(accToken)
	db := database.DbConn()
	defer db.Close()
	selDB, err := db.Query("SELECT username from refresh_token where access_token=? and refresh_token=? and is_active='1'", accToken, refToken)
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	defer selDB.Close()
	if selDB.Next() == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} // else {

	err = selDB.Scan(&usr)
	if err != nil {
		WriteLogFile(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	updDB, err := db.Prepare("UPDATE refresh_token SET is_active=? WHERE username=? and refresh_token=?")
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	updDB.Exec(0, usr, refToken)
	defer updDB.Close()
	//}

	// Now, create a new token for the current use, with a renewed expiration time
	// expirationTime := time.Now().Add(3600 * time.Second)
	// claims.ExpiresAt = expirationTime.Unix()
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenStringRefresh, err := token.SignedString(jwtKey)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	expirationTime := time.Now().Add(3600 * time.Second).Unix()
	claims = &model.Claims{
		Username: usr,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenStringRefresh, err := token.SignedString(jwtKey)
	if err != nil {
		WriteLogFile(err)
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ_-&" +
		"abcdefghijklmnopqrstuvwxyz%*" +
		"0123456789")
	length := 20
	var bRefresh strings.Builder
	for i := 0; i < length; i++ {
		bRefresh.WriteRune(chars[rand.Intn(len(chars))])
	}
	var jwtTokenRefresh model.JwtToken
	jwtTokenRefresh.AccessToken = tokenStringRefresh
	jwtTokenRefresh.TokenType = "bearer"
	jwtTokenRefresh.Expiry = "3600"
	jwtTokenRefresh.RefreshToken = bRefresh.String()
	json.NewEncoder(w).Encode(jwtTokenRefresh)
	updatedAt := time.Now()
	insForm, err := db.Prepare("Insert into token(username,access_token,expiration,role,created_at) values (?,?,?,?,?)")
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	defer insForm.Close()
	fmt.Println(tokenStringRefresh)
	insForm.Exec(usr, tokenStringRefresh, expirationTime, Role, updatedAt.Format("2006-01-02T15:04:05"))
	query := "Insert into refresh_token(username,access_token,refresh_token,created_at) values (?,?,?,?)"
	insert2, err := db.Prepare(query)
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	insert2.Exec(usr, tokenStringRefresh, bRefresh.String(), updatedAt.Format("2006-01-02T15:04:05"))
	defer insert2.Close()

}
