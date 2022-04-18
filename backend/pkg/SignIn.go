package social

import (
	"context"
	"fmt"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// connect db globally so all funcs can use client rather than waste connections
var clientOptions = options.Client().ApplyURI("mongodb+srv://andru:1q1q1q@cluster0.tccti.mongodb.net/cluster0?retryWrites=true&w=majority") // Set client options

var  client, err = mongo.Connect(context.TODO(), clientOptions) // Connect to MongoDB
   
var err1 = client.Ping(context.TODO(), nil) // Check the connection

var sess, err2 = session.NewSession(&aws.Config{ //start a aws session by setting the region
Region: aws.String("us-east-2")},
)

func Signin (c *gin.Context) {// takes id and sets up bucket and mongodb
    userid := c.Param("userid") // get id from url request
    
    type usrsignin struct { 
    Username     string  `bson:"Username" json:"Username"`
    Password  string  `bson:"Password" json:"Password"`
    }
    
    var reqbody usrsignin // declare new instance of struct type
    
    if err := c.BindJSON(&reqbody); err != nil {
        fmt.Println(err)
        return
    }
 
    
    collection := client.Database("datingapp").Collection("userdata")// connect to db and collection
    
    result := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request.
	// call the collection's Find() method and return object into result
    err = collection.FindOne(ctx, bson.M{"Username": userid}).Decode(&result)

    if result.Password == reqbody.Password {
        fmt.Println("password matches")
        token := makejwt(userid, true) // make jwt
        c.JSON(http.StatusOK, gin.H{ 
				"token": token,
			     })
    } else {
        fmt.Println("username or password is not a match")
        c.IndentedJSON(http.StatusUnauthorized, nil)
    }
   
}