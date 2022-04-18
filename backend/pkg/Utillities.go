package social

import (
	"context"
	"fmt"
    "log"
    "strings"
    "os"
    "encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
    "github.com/gin-gonic/gin"
)

// connect db globally so all funcs can use client rather than waste connections
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


func connectedmngo () { // prints connected if all error checks passed
    if err != nil || err1 != nil {
        log.Fatal(err)
    }else {
    fmt.Println("Connected to MongoDB!") 
    }
}

func connectedaws () { // prints connected to aws if all error checks passed
    if err2 != nil {
        log.Fatal(err)
    }else {
    fmt.Println("Connected to MongoDB!") 
    }
}


func listbuckets () {
	// Create S3 service client
   svc := s3.New(sess)
   
   // list all buckets
   result, err := svc.ListBuckets(nil)
   if err != nil {
	   fmt.Println(err)
   }

   fmt.Println("Buckets:")

   for _, b := range result.Buckets { // loop through results print name and date created
	 fmt.Printf("* %s created on %s\n",
	  aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
   }
}


func createbucket (bucketname string) {// creates a s3 bucket with the name passed to it
	// Create S3 service client
   svc := s3.New(sess)
   
   _, err = svc.CreateBucket(&s3.CreateBucketInput{
   Bucket: aws.String(bucketname+uniqueadr),// make bucket name unique
   })
   if err != nil {
		fmt.Println("Unable to create bucket %q, %v", bucketname+uniqueadr, err)
   }

   // Wait until bucket is created before finishing
   fmt.Printf("Waiting for bucket %q to be created...\n", bucketname+uniqueadr)

   err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
   Bucket: aws.String(bucketname+uniqueadr),
   })
   
   if err != nil { // check bucket is created 
	   fmt.Println("Error occurred while waiting for bucket to be created, %v", bucketname+uniqueadr)
   }else {
	   fmt.Printf("Bucket %q successfully created\n", bucketname+uniqueadr)
   }
   
   publicbucket(bucketname) // make bucket public read only
}


func publicbucket (bucket string) {// make bucket public and read only
	// Create S3 service client
   svc := s3.New(sess)
   
   
   readOnlyAnonUserPolicy := map[string]interface{}{ // add policy to map so can be sent
	   "Version": "2012-10-17",
	   "Statement": []map[string]interface{}{
		   {
			   "Sid":       "AddPerm",
			   "Effect":    "Allow",
			   "Principal": "*",
			   "Action": []string{
				   "s3:GetObject",
			   },
			   "Resource": []string{
				   fmt.Sprintf("arn:aws:s3:::%s/*", bucket+uniqueadr),
			   },
		   },
	   },
   }
   
   
   policy, err := json.Marshal(readOnlyAnonUserPolicy)
   
   if err != nil {
	   fmt.Println("Unable to marshal json %v", err)
   }
   
   _, err = svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
	   Bucket: aws.String(bucket+uniqueadr),
	   Policy: aws.String(string(policy)),
   })
   
   if err != nil {
	   fmt.Println("Unable to update bucket policy %v", err)
   } else {
	   fmt.Printf("Successfully set bucket %q's policy\n", bucket+uniqueadr)
   }

}


func listitems (bucket string) {
	// Create S3 service client
   svc := s3.New(sess)
   
   resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
   if err != nil {
	   fmt.Println("Unable to list items in bucket %q, %v", bucket, err)
   }

   for _, item := range resp.Contents { // loop through bucket contents and use attributes to print info
	   fmt.Println("Name:         ", *item.Key)
	   fmt.Println("Last modified:", *item.LastModified)
	   fmt.Println("Size:         ", *item.Size)
	   fmt.Println("Storage class:", *item.StorageClass)
	   fmt.Println("")
   }

}

func GetTempLoc(filename string) string { // get the temp location of files sent in post requests to be used for aws upload without having to store data locally
    return strings.TrimRight(os.TempDir(), "/") + "/" + filename
}


func getfiletype(filename string) string { // takes filname and returns filetype

	runed:=[]rune(filename)
	var result [][]rune
	
    for i:= len(runed)-1 ; i>=0; i-- { // looks for . and cuts the file type to new array
		if runed[i] == 46 {
		  result = append(result, runed[i:])
		  break
		}
	}

	return(string(result[0]))
}

func exitErrorf(msg string, args ...interface{}) {// hadles aws errors
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}


func CORSMiddleware() gin.HandlerFunc { // cors func to allow body and acces from anywhere
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}