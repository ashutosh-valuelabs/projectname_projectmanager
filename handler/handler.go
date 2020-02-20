package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Commander : structure for commander
type Commander struct{}

// type TokenData struct {
// 	UserName string
// 	Role     string
// }

// UserName : for user authentication
var UserName, Role string
var expiration int64

// type key int

// var id = key(1)

// func Set(ctx context.Context) context.Context {
// 	db := database.DbConn()
// 	user, err := db.Query("SELECT username FROM token WHERE access_token=? AND is_active = '1'", reqToken)

// 	if err != nil {
// 		WriteLogFile(err)
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 	}
// 	defer user.Close()

// 	return context.WithValue(ctx, id, "secret")
// }

// func Get(ctx context.Context) (string, bool) {
// 	val, ok := ctx.Value(id).(string)
// 	return val, ok
// }

//authMiddleware:Auth Middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//ctx := newContextWithRequestID(r.Context(), r)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		var error model.Error
		db := database.DbConn()
		defer db.Close()
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		user, err := db.Query("SELECT username,role, expiration FROM token WHERE access_token=? AND is_active = '1'", reqToken)
		defer user.Close()
		if err != nil {
			WriteLogFile(err)
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		if user.Next() != false {
			err := user.Scan(&UserName, &Role, &expiration)
			if err != nil {
				WriteLogFile(err)
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
			Role = strings.ToLower(Role)
			tm := time.Unix(expiration, 0)
			currentTime := time.Now()
			fmt.Println(currentTime, tm)
			inTime := currentTime.Before(tm)
			if inTime == false {
				updDB, err := db.Prepare("UPDATE token SET is_active=? WHERE username=? and access_token=?")
				if err != nil {
					WriteLogFile(err)
					panic(err.Error())
				}
				updDB.Exec(0, UserName, reqToken)
				defer updDB.Close()
				// updDB, err = db.Prepare("UPDATE refresh_token SET is_active=? WHERE username=? and access_token=?")
				// if err != nil {
				// 	WriteLogFile(err)
				// 	panic(err.Error())
				// }
				// updDB.Exec(0, UserName, reqToken)
				// defer updDB.Close()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				error.Message = fmt.Sprintf("Error validating access token: Session has expired on %s", tm)
				json.NewEncoder(w).Encode(error)
				return
			}
			// user, err = db.Query("SELECT role FROM login WHERE username=?", UserName)
			// defer user.Close()
			// if err != nil {
			// 	WriteLogFile(err)
			// 	fmt.Println(err)
			// 	w.WriteHeader(http.StatusBadRequest)
			// }
			// if user.Next() != false {
			// 	err := user.Scan(&Role)
			// 	if err != nil {
			// 		WriteLogFile(err)
			// 		fmt.Println(err)
			// 		w.WriteHeader(http.StatusBadRequest)
			// 	}
			// }
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

//loggingMiddleware:Logging Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

//HandleRequests : handler function
func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	secure := myRouter.PathPrefix("").Subrouter()
	secure.Use(authMiddleware)
	rout := getconfig()
	c := &Commander{}
	for i := 0; i < len(rout.R); i++ {
		fmt.Println(rout.R[i].Path, rout.R[i].Callback, rout.R[i].Method, rout.R[i].Authorization)
		m := reflect.ValueOf(c).MethodByName(rout.R[i].Callback)
		Call := m.Interface().(func(http.ResponseWriter, *http.Request))

		if rout.R[i].Authorization == "YES" {

			secure.HandleFunc(rout.R[i].Path, Call).Methods(rout.R[i].Method)
		} else {
			myRouter.HandleFunc(rout.R[i].Path, Call).Methods(rout.R[i].Method)
		}

	}
	myRouter.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":8008", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(myRouter)))
}

func getconfig() (c model.Conf) {
	yamlFile, err := ioutil.ReadFile("routes.yaml")
	if err != nil {
		WriteLogFile(err)
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), &c)
	//fmt.Println("1")
	if err != nil {
		WriteLogFile(err)
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

//SetupResponse : to setup access control on requests
func SetupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*req).Method == "OPTIONS" {
		(*w).Header().Set("Access-Control-Max-Age", "86400")
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

//WriteLogFile : error logging
func WriteLogFile(err error) {
	f := OpenLogFile()
	pc, fn, line, _ := runtime.Caller(1)
	fmt.Println(err)
	fmt.Println(pc, fn, line)
	log.SetOutput(f)
	log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	defer f.Close()
}

// OpenLogFile : this function will open log file and return the file writer
func OpenLogFile() (f *os.File) {
	f, err := os.OpenFile("logs/output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	return f
}

//BadRequest : to handle bad requests
func BadRequest(w http.ResponseWriter, err error) {
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
