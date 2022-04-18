package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    // "github.com/andru100/Social-Network/backend/pkg"
    "Social"

)


func main() {
        

        // connectedmngo(err, err1)// chk errors then display "Connection succesful"
        // connectedaws(err2)// chk errors then display "Connection succesful"
        router := gin.New()  
        router.Use(social.CORSMiddleware())
        router.POST("/getCmt", social.Stdget) // gets comments for user id parsed in url
        router.PUT("/postMsg", social.PostMsg) 
        router.PUT("/updatebio", social.Updatebio)// add user bio
        router.POST("/postfile/:userid", social.Postfile)// posts profile pic and users media
        router.POST("/signup/:userid", social.Signup) // used to take user id on sign up and create s3bucket and mongo doc
        router.POST("/signin/:userid", social.Signin) // signs user in    
        router.POST("/chkauth", social.Chkauth) // checks for authentication using jwt
        router.Run(":4001")
        //router.RunTLS(":4001", "./server.pem", "./server.key")
}