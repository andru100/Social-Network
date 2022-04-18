package social

import (
	"context"
	"fmt"
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

func Updatebio(c *gin.Context) {// updates user bio section


    type stdputstruct struct { // create struct with data types of incoming json object- used to find out what type of query is needed
        Username     string  `bson:"Username" json:"Username"`
        Updatetype  string  `bson:"Updatetype" json:"Updatetype"`
        Value2updt  string `bson:"Value2updt" json:"Value2updt"`
        Key2updt string     `bson:"Key2updt" json:"Key2updt"`
    }
    
    var reqbody stdputstruct
   
    
    if err := c.BindJSON(&reqbody); err != nil {
        fmt.Println(err)
        return
    }
 
    collection := Client.Database("datingapp").Collection("userdata")

    filter := bson.M{"Username": reqbody.Username}

    update := bson.D{
        {reqbody.Updatetype, bson.D{
            {reqbody.Key2updt, reqbody.Value2updt},
        }},
    }
    
    //put to db
    _, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request. call the collection's Find() method and return object into result
    err = collection.FindOne(ctx, bson.M{"Username": reqbody.Username}).Decode(&currentDoc)
    
    c.IndentedJSON(http.StatusOK, currentDoc)
}
