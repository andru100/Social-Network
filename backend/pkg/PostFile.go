package social

import (
	"context"
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
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


func Postfile(c *gin.Context) {// post file takes file from request form, runs upload func, puts in s3, returns s3 address.

    type formData struct {
        User     string `form:"user" binding:"required"`
        Type string `form:"type" binding:"required"`
    }


    file, header, err := c.Request.FormFile("file") // get file from request body
    if err != nil {
        fmt.Println(err)
    }
 
   filename := header.Filename
    
   fileread, err := ioutil.ReadAll(file) // read the file to variable 

    userid := c.Param("userid") // get id from url request
    requestType := c.PostForm("type")

    collection := client.Database("datingapp").Collection("userdata")

    imgaddress := uploaditem(userid, filename, fileread)// call upload func returns uploaded img url
    
    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request.
    err = collection.FindOne(ctx, bson.M{"Username": userid}).Decode(&currentDoc)
    //create filter 
    filter := bson.M{"Username": userid}
    update := bson.D{}
    if requestType == "profPic" {
        currentDoc.Profpic= imgaddress //replace url to profpic section of user object
        update = bson.D{
            {"$set", bson.D{
                {"Profpic", currentDoc.Profpic},
            }},
        }
    } else if requestType == "addPhotos" {
        currentDoc.Photos= append(currentDoc.Photos, imgaddress) //append to list of users photo urls
        update = bson.D{
            {"$set", bson.D{
                {"Photos", currentDoc.Photos},
            }},
        }
    }
    
    //put to db
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    c.IndentedJSON(http.StatusOK, currentDoc)
}