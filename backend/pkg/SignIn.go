package social

import (
	"context"
	"fmt"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
 
    
    collection := Client.Database("datingapp").Collection("userdata")// connect to db and collection
    
    result := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request.
	// call the collection's Find() method and return object into result
    err := collection.FindOne(ctx, bson.M{"Username": userid}).Decode(&result)
    CheckError(err)

    if result.Password == reqbody.Password {
        fmt.Println("password matches")
        token := Makejwt(userid, true) // make jwt
        c.JSON(http.StatusOK, gin.H{ 
				"token": token,
			     })
    } else {
        fmt.Println("username or password is not a match")
        c.IndentedJSON(http.StatusUnauthorized, nil)
    }
   
}