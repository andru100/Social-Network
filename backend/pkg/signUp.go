package social

	import (
		"context"
		"fmt"
		"log"
		"time"
		"net/http"
		"github.com/gin-gonic/gin"
		"go.mongodb.org/mongo-driver/bson/primitive"
		"go.mongodb.org/mongo-driver/bson"
	)

func Signup (c *gin.Context) {// takes id and sets up bucket and mongodb
    userid := c.Param("userid") // get id from url request
    fmt.Println("userid is ", userid)
    Createbucket(userid) // create bucket to store users files

    var reqbody Usrsignin // declare new instance of struct type

    if err := c.BindJSON(&reqbody); err != nil {
        fmt.Println(err)
        return
    }
 
    
    collection := Client.Database("datingapp").Collection("userdata")// connect to db and collection.
    
    //post to db
    insertResult, err := collection.InsertOne(context.TODO(), reqbody)
    if err != nil {
        log.Fatal(err)
    }
    
    // Declare a struct to create an “empty” MongoDB document that can be used to store values returned by the API call
    result := MongoFields{}

    // Declare Context type object for managing multiple API requests 
    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request.
    idstr := insertResult.InsertedID
    err = collection.FindOne(ctx, bson.M{"_id": idstr}).Decode(&result)

    //----------------------------------sign in the user return jwt----------------------------------------------//
    
    //create struct to hold retrived mongo doc to check password
    type passwordChk struct {
    Key string `json:"key,omitempty"`
    ID primitive.ObjectID `bson:"_id, omitempty"` 
    Username     string  `bson:"Username" json:"Username"`
    Password  string  `bson:"Password" json:"Password"`
    }
    
    result1 := passwordChk{}

    err = collection.FindOne(ctx, bson.M{"Username": reqbody.Username}).Decode(&result1)
    
    if result1.Password == reqbody.Password {
        fmt.Println("password matches")
        token := Makejwt(userid, true) // make jwt with user id and auth true
        c.JSON(http.StatusOK, gin.H{ //make header with token in and send
				"token": token,
		})
    } else {
        fmt.Println("username or password is not a match")
        c.IndentedJSON(http.StatusUnauthorized, nil)
    }
   
}

