package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

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

func main() {
        

        // connectedmngo(err, err1)// chk errors then display "Connection succesful"
        // connectedaws(err2)// chk errors then display "Connection succesful"
        router := gin.New()  
        router.Use(CORSMiddleware())
        router.POST("/getCmt", stdget) // gets comments for user id parsed in url
        router.PUT("/postMsg", postMsg) 
        router.PUT("/updatebio", updatebio)// add user bio
        router.POST("/postfile/:userid", postfile)// posts profile pic and users media
        router.POST("/signup/:userid", signup) // used to take user id on sign up and create s3bucket and mongo doc
        router.POST("/signin/:userid", signin) // signs user in    
        router.POST("/chkauth", chkauth) // checks for authentication using jwt
        router.Run(":4001")
        //router.RunTLS(":4001", "./server.pem", "./server.key")
}