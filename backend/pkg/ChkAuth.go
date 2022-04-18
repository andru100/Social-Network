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
)

/*---------------------------------------------------------------------------------------------*/

// connect db globally so all funcs can use client rather than waste connections
var clientOptions = options.Client().ApplyURI("mongodb+srv://andru:1q1q1q@cluster0.tccti.mongodb.net/cluster0?retryWrites=true&w=majority") // Set client options

var  client, err = mongo.Connect(context.TODO(), clientOptions) // Connect to MongoDB
   
var err1 = client.Ping(context.TODO(), nil) // Check the connection

var sess, err2 = session.NewSession(&aws.Config{ //start a aws session by setting the region
Region: aws.String("us-east-2")},
)
var uniqueadr = "ajh46unique"

type msgCmts struct {
Username string `bson:"Username" json:"Username"`
Comment string `bson:"Comment" json:"Comment"`
Profpic string `bson:"Profpic" json:"Profpic"`
}

type Likes struct {
Username string `bson:"Username" json:"Username"`
Profpic string `bson:"Profpic" json:"Profpic"`
}

type PostData struct {
Username     string  `bson:"Username" json:"Username"`
SessionUser     string  `bson:"SessionUser" json:"SessionUser"`
MainCmt string  `bson:"MainCmt" json:"MainCmt"`
PostNum int  `bson:"PostNum" json:"PostNum"`
Time string  `bson:"Time" json:"Time"`
TimeStamp  int64  `bson:"TimeStamp" json:"TimeStamp"`
Date string  `bson:"Date" json:"Date"`
Comments [] msgCmts  `bson:"Comments" json:"Comments"`
Likes [] Likes `bson:"Likes" json:"Likes"`
}

type usrsignin struct { 
Username     string  `bson:"Username" json:"Username"`
Password  string  `bson:"Password" json:"Password"`
Email  string  `bson:"Email" json:"Email"`
Bio string `bson:"Bio" json:"Bio"`
Photos [] string `bson:"Photos" json:"Photos"`
LastCommentNum int  `bson:"LastCommentNum" json:"LastCommentNum"`
LikeSent Likes   `bson:"LikeSent" json:"LikeSent"`
Posts  []PostData  `bson:"Posts" json:"Posts"`
}

//struct to hold retrived mongo doc
type MongoFields struct {
Key string `json:"key,omitempty"`
ID primitive.ObjectID `bson:"_id, omitempty"`  
Username     string  `bson:"Username" json:"Username"`
Password  string  `bson:"Password" json:"Password"`
Email  string  `bson:"Email" json:"Email"`
Bio string `bson:"Bio" json:"Bio"`
Profpic string `bson:"Profpic" json:"Profpic"`
Photos [] string `bson:"Photos" json:"Photos"`
LastCommentNum int  `bson:"LastCommentNum" json:"LastCommentNum"`
LikeSent Likes   `bson:"LikeSent" json:"LikeSent"`
Posts  []PostData  `bson:"Posts" json:"Posts"`
}

func chkauth (c *gin.Context) { // checks for authentication

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

