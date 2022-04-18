package main

import (
        "github.com/andru100/Social-Network/backend/pkg"
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