package main

import (
    "github.com/gin-gonic/gin"
)


func main() {
        // connectedmngo(err, err1)// chk errors then display "Connection succesful"
        // connectedaws(err2)// chk errors then display "Connection succesful"
        router := gin.New()  
        router.Use(CORSMiddleware())
        router.POST("/getCmt", Stdget) // gets comments for user id parsed in url
        router.PUT("/postMsg", PostMsg) 
        router.PUT("/updatebio", Updatebio)// add user bio
        router.POST("/postfile/:userid", Postfile)// posts profile pic and users media
        router.POST("/signup/:userid", Signup) // used to take user id on sign up and create s3bucket and mongo doc
        router.POST("/signin/:userid", Signin) // signs user in    
        router.POST("/chkauth", Chkauth) // checks for authentication using jwt
        router.Run(":4001")
        //router.RunTLS(":4001", "./server.pem", "./server.key")
}