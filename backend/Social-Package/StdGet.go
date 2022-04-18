package social

import (
	"context"
	"fmt"
	"sort"
	"log"
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

func stdget(c *gin.Context) {// gets comments for a specified user/ all users if on home feed page
    
    type qrystruct struct { 
      Page     string  `bson:"Page" json:"Page"`
      UserName string   `bson:"UserName" json:"UserName"`
    }
    
    var qry qrystruct 
    if err := c.Bind(&qry); err != nil {
        fmt.Println(err)
        return
    }
 
    
    collection := client.Database("datingapp").Collection("userdata")// connect to db and collection.

    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

    var allPosts []PostData
    
    if qry.Page == "all" { // used when home feed wants all comments from all users
        findOptions := options.Find()
        findOptions.SetLimit(2)

        var results []*MongoFields

        // Passing bson.D{{}} as the filter matches all documents in the collection
        cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
        if err != nil {
            log.Fatal(err)
        }

        // Iterating through cursor decode documents one at a time
        for cur.Next(context.TODO()) {
            
            var elem MongoFields
            err := cur.Decode(&elem)
            if err != nil {
                log.Fatal(err)
            }

            results = append(results, &elem)
        }

        if err := cur.Err(); err != nil {
            log.Fatal(err)
        }

        cur.Close(context.TODO())

        for _, record := range results {
            for _, posts := range record.Posts{
              allPosts = append(allPosts, posts)
            }
        }

        sort.Slice(allPosts, func(i, j int) bool { //Sort posts or comments by time descending
            return allPosts[i].TimeStamp > allPosts[j].TimeStamp
        })

        var json2send MongoFields
        json2send.Posts = allPosts // send in this format so front end can use same func to access single user profile requests
        err = collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)
        json2send.Profpic = currentDoc.Profpic
        json2send.Bio = currentDoc.Bio 
        json2send.Photos = currentDoc.Photos 
        
        c.IndentedJSON(http.StatusOK, json2send)
    } else if qry.Page == "media" { // if page is users media section
        err = collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)

        c.IndentedJSON(http.StatusOK, currentDoc)

    } else{

        err = collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)

        sort.Slice(currentDoc.Posts, func(i, j int) bool {
            return currentDoc.Posts[i].TimeStamp > currentDoc.Posts[j].TimeStamp
        })

        c.IndentedJSON(http.StatusOK, currentDoc)
     }
   
}