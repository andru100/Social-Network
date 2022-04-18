package social

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/dgrijalva/jwt-go"
)

func Chkauth (c *gin.Context) { // checks for authentication

    var jwtKey = []byte("AllYourBase")
    
    type Claims struct {
        Username string `json:"username"`
        jwt.StandardClaims
    }
    
    
    type jwtdata struct { // for user sign in  data
    Data1     string  `bson:"Data1" json:"Data1"`
    }
    
    var reqbody jwtdata 
    
    if err := c.BindJSON(&reqbody); err != nil { 
        fmt.Println(err)
        return
    }
    
    tknStr := reqbody.Data1           
    
    // check token is valid return username in response
	claims := &Claims{}

	// Parse the JWT string and store the result in claims.
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.IndentedJSON(http.StatusUnauthorized, nil)
			return
		}
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}
	if !tkn.Valid {
		c.IndentedJSON(http.StatusUnauthorized, nil)
		return
	}

    type authd struct {
        AuthdUser     string  `bson:"AuthdUser" json:"AuthdUser"`
    }

    var auth authd 
    
    auth.AuthdUser = claims.Username 

	// Finally, return the welcome message to the user, along with their username given in the token
    c.IndentedJSON(http.StatusOK, auth)
}

