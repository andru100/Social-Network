package social 

import (
	"context"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// connect db globally so all funcs can use client rather than waste connections
var clientOptions = options.Client().ApplyURI("mongodb+srv://andru:1q1q1q@cluster0.tccti.mongodb.net/cluster0?retryWrites=true&w=majority") // Set client options

var  client, err = mongo.Connect(context.TODO(), clientOptions) // Connect to MongoDB
   
var err1 = client.Ping(context.TODO(), nil) // Check the connection

var sess, err2 = session.NewSession(&aws.Config{ //start a aws session by setting the region
Region: aws.String("us-east-2")},
)
var uniqueadr = "ajh46unique"

func Uploaditem (bucket string, filename string, filebytes []byte) string {// upload file to s3 with the bucket name and file adress passed to it

    tmpfile, err := ioutil.TempFile("", "example")// create temp file using naming convention.. it'll ad random stuff
    // empty string in first arg tells it to go to default temp dir set by os
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(filebytes); err != nil { //write file from []bytes given by io readall
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
    
    
    file, err := os.Open(tmpfile.Name())// open file using the temp dir and temp name created.
    if err != nil {
        fmt.Println("Unable to open file %q, %v", err)
    }

    defer file.Close()// clean up
 
    uploader := s3manager.NewUploader(sess)
    
    result, err := uploader.Upload(&s3manager.UploadInput{// upload file
    Bucket: aws.String(bucket+uniqueadr),
    Key: aws.String(filename),
    Body: file,
    })
    
    if err != nil {
        // Print the error and exit.
        fmt.Printf("Unable to upload %q to %q, %v\n", filename, bucket+uniqueadr, err)
    } else {
        fmt.Println("return result after upload is", result)
    }

    fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket+uniqueadr)
    return result.Location

} 