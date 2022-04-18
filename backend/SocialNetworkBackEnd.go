package main

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
        "github.com/aws/aws-sdk-go/service/s3"
        "github.com/aws/aws-sdk-go/service/s3/s3manager"
        "strings"
        "os"
        "io/ioutil"
        "encoding/json"
        "github.com/dgrijalva/jwt-go"
)

func main() {
        connectedmngo()// chk errors then display "Connection succesful"
        connectedaws()// chk errors then display "Connection succesful"
        router := gin.New()  
        router.Use(CORSMiddleware())
        router.POST("/getCmt", stdget) // gets comments for user id parsed in url
        router.PUT("/postMsg", postMsg) 
        router.PUT("/updatebio", updatebio)// add user bio
        router.POST("/postfile/:userid", postfile)// posts profile pic and users media
        router.POST("/signup/:userid", signup) // used to take user id on sign up and create s3bucket and mongo doc
        router.POST("/signin/:userid", signin) // signs user in    
        router.POST("/chkauth", chkauth) // checks for authentication using jwt
        listbuckets() // list current buckets on startup
        router.Run(":4001")
        //router.RunTLS(":4001", "./server.pem", "./server.key")
}