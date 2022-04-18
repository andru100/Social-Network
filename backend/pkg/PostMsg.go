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

// connect db globally so all funcs can use client rather than waste connections
var clientOptions = options.Client().ApplyURI("mongodb+srv://andru:1q1q1q@cluster0.tccti.mongodb.net/cluster0?retryWrites=true&w=majority") // Set client options

var  client, err = mongo.Connect(context.TODO(), clientOptions) // Connect to MongoDB
   
var err1 = client.Ping(context.TODO(), nil) // Check the connection

var sess, err2 = session.NewSession(&aws.Config{ //start a aws session by setting the region
Region: aws.String("us-east-2")},
)

func PostMsg(c *gin.Context) {// std post creates doc from query, finds it and returns it


    type stdputstruct struct { // create struct with data types of incoming json object- used to find out what type of query is needed
        Username     string  `bson:"Username" json:"Username"`
        Updatetype  string  `bson:"Updatetype" json:"Updatetype"`
        IsReply  string  `bson:"IsReply" json:"IsReply"`// see if msg is a reply to post
        LikeSent Likes   `bson:"LikeSent" json:"LikeSent"`
        ReplyCmt msgCmts `bson:"ReplyCmt" json:"ReplyCmt"`
        PostIndx   int     `bson:"PostIndx" json:"PostIndx"` // used to add comments to specific post
        Key2updt string     `bson:"Key2updt" json:"Key2updt"`
        Value2updt  PostData `bson:"Value2updt" json:"Value2updt"`
    }
    
    var reqbody stdputstruct
   
    
    if err := c.BindJSON(&reqbody); err != nil {
        fmt.Println(err)
        return
    }
    
    collection := client.Database("datingapp").Collection("userdata")
    
    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request. call the collection's Find() method and return object into result
    err = collection.FindOne(ctx, bson.M{"Username": reqbody.Value2updt.Username}).Decode(&currentDoc)

    //create filter 
    filter := bson.M{"Username": ""}

    // find out if msg to be added was a reply to a comment, a new post or a like.
    if reqbody.IsReply == "isResponse" {
        
        for i:=0; i<len(currentDoc.Posts); i++ {
            if currentDoc.Posts[i].PostNum == reqbody.PostIndx {
                currentDoc.Posts[i].Comments= append(currentDoc.Posts[i].Comments, reqbody.ReplyCmt) // add reply to post 
                filter = bson.M{"Username": currentDoc.Posts[i].Username} // change filter to post maker
            }
        }
    } else if reqbody.IsReply == "isCmt" {
        currentDoc.LastCommentNum += 1 
        reqbody.Value2updt.PostNum = currentDoc.LastCommentNum
        currentDoc.Posts= append(currentDoc.Posts, reqbody.Value2updt) 
        filter = bson.M{"Username": reqbody.Value2updt.Username} // set filter to username of post creator
        
    } else if reqbody.IsReply == "cmtLiked" { 
        for i:=0; i<len(currentDoc.Posts); i++ {
            if currentDoc.Posts[i].PostNum == reqbody.PostIndx {
                currentDoc.Posts[i].Likes = append(currentDoc.Posts[i].Likes, reqbody.LikeSent)  // add like to post 
                filter = bson.M{"Username": currentDoc.Posts[i].Username} // change filter to post maker
            }
        }
    } 

    update := bson.D{
        {reqbody.Updatetype, bson.D{
            {reqbody.Key2updt, currentDoc.Posts},
        }},
    }
    
    //put to db
    _, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    //update post index count

    update = bson.D{
        {reqbody.Updatetype, bson.D{
            {"LastCommentNum", currentDoc.LastCommentNum},
        }},
    }
    
    //put to db
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    result := MongoFields{}

    ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)

    var allPosts []PostData
    
    //check if request is made from profile page or news feed
    if reqbody.Username == "all" {

        findOptions := options.Find()
        findOptions.SetLimit(2)

        var results []*MongoFields

        cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
        if err != nil {
            log.Fatal(err)
        }

        //Iterate through the cursor decode documents one at a time
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

        sort.Slice(allPosts, func(i, j int) bool {
            return allPosts[i].TimeStamp > allPosts[j].TimeStamp
        })

        var json2send MongoFields
        json2send.Posts = allPosts // send in this format so front end can use same function to process other requests
        err = collection.FindOne(ctx, bson.M{"Username": reqbody.Value2updt.SessionUser}).Decode(&currentDoc)
        json2send.Profpic = currentDoc.Profpic
        json2send.Bio = currentDoc.Bio 
        json2send.Photos = currentDoc.Photos 
        c.IndentedJSON(http.StatusOK, json2send)
    } else {
        err = collection.FindOne(ctx, bson.M{"Username": reqbody.Value2updt.Username}).Decode(&result)
        sort.Slice(result.Posts, func(i, j int) bool {
            return result.Posts[i].TimeStamp > result.Posts[j].TimeStamp
        })

        c.IndentedJSON(http.StatusOK, result)
    }
   
}