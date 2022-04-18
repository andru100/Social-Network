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

	// connect to db 
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


func signup (c *gin.Context) {// takes id and sets up bucket and mongodb
    userid := c.Param("userid") // get id from url request
    fmt.Println("userid is ", userid)
    createbucket(userid) // create bucket to store users files

    var reqbody usrsignin // declare new instance of struct type

    if err := c.BindJSON(&reqbody); err != nil {
        fmt.Println(err)
        return
    }
 
    
    collection := client.Database("datingapp").Collection("userdata")// connect to db and collection.
    
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
        token := makejwt(userid, true) // make jwt with user id and auth true
        c.JSON(http.StatusOK, gin.H{ //make header with token in and send
				"token": token,
		})
    } else {
        fmt.Println("username or password is not a match")
        c.IndentedJSON(http.StatusUnauthorized, nil)
    }
   
}

