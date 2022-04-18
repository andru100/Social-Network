package social

import (
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func Makejwt (userid string, isauth bool) string {
	mySigningKey := []byte("AllYourBase")//base for ecoding eg private key
 
	type Claims struct { // for claim settings
	 Username string `json:"username"`
	 jwt.StandardClaims
	}
 
	 // Create the Claims
	// Declare the expiration time of the token
	 expirationTime := time.Now().Add(5 * time.Minute)
	 // Create the JWT claims, which includes the username and expiry time
	 claims := &Claims{
		 Username: userid,
		 StandardClaims: jwt.StandardClaims{
			 // In JWT, the expiry time is expressed as unix milliseconds
			 ExpiresAt: expirationTime.Unix(),
		 },
	 }
	 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)// converts to json and signs/encodes
	 ss, err := token.SignedString(mySigningKey)// encode to string with chosen base string
	 fmt.Printf("jwt created %v %v", ss, err)
	 return ss
 
 
 }