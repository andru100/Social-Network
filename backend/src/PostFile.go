package social

import (
	"context"
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

    collection := Client.Database("datingapp").Collection("userdata")

    imgaddress := Uploaditem(userid, filename, fileread)// call upload func returns uploaded img url
    
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