package model

import "github.com/dgrijalva/jwt-go"

// Credentials : user credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Claims : jwtToken
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JwtToken : access token
type JwtToken struct {
	AccessToken  string `json:"token"`
	TokenType    string `json:"token_type"`
	Expiry       string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// Client : LDAP functions
type Client interface {
	Auth(username, password string) error
	//Role(username string) error
}

// Config : LDAP configurations
type Config struct {
	BaseDN string // base directory, ex. "CN=Users,DC=Company"
	ROUser User   // the read-only user for initial bind
	Host   string // the ldap host and port, ex. "ldap.directory.com:389"
	Filter string // defaults to "sAMAccountName" for AD
	Title  string
}

// User holds the name and pass required for initial read-only bind.
type User struct {
	Name     string
	Password string
}
