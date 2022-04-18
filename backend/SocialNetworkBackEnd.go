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

/*---------------------------------------------------------------------------------------------*/

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

/*---------------------------------------------------------------------------------------------*/

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

/*---------------------------------------------------------------------------------------------*/

func connectedmngo () { // prints connected if all error checks passed
    if err != nil || err1 != nil {
        log.Fatal(err)
    }else {
    fmt.Println("Connected to MongoDB!") 
    }
}

/*---------------------------------------------------------------------------------------------*/

func connectedaws () { // prints connected to aws if all error checks passed
    if err2 != nil {
        log.Fatal(err)
    }else {
    fmt.Println("Connected to MongoDB!") 
    }
}

/* -------------------------------------------------------------------------------------------------------------------------*/


// func signup (c *gin.Context) {// takes id and sets up bucket and mongodb
//     userid := c.Param("userid") // get id from url request
//     fmt.Println("userid is ", userid)
//     createbucket(userid) // create bucket to store users files

//     var reqbody usrsignin // declare new instance of struct type

//     if err := c.BindJSON(&reqbody); err != nil {
//         fmt.Println(err)
//         return
//     }
 
    
//     collection := client.Database("datingapp").Collection("userdata")// connect to db and collection.
    
//     //post to db
//     insertResult, err := collection.InsertOne(context.TODO(), reqbody)
//     if err != nil {
//         log.Fatal(err)
//     }
    
//     // Declare a struct to create an “empty” MongoDB document that can be used to store values returned by the API call
//     result := MongoFields{}

//     // Declare Context type object for managing multiple API requests 
//     ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
//     // Find the document that mathes the id from the request.
//     idstr := insertResult.InsertedID
//     err = collection.FindOne(ctx, bson.M{"_id": idstr}).Decode(&result)

//     //----------------------------------sign in the user return jwt----------------------------------------------//
    
//     //create struct to hold retrived mongo doc to check password
//     type passwordChk struct {
//     Key string `json:"key,omitempty"`
//     ID primitive.ObjectID `bson:"_id, omitempty"` 
//     Username     string  `bson:"Username" json:"Username"`
//     Password  string  `bson:"Password" json:"Password"`
//     }
    
//     result1 := passwordChk{}

//     err = collection.FindOne(ctx, bson.M{"Username": reqbody.Username}).Decode(&result1)
    
//     if result1.Password == reqbody.Password {
//         fmt.Println("password matches")
//         token := makejwt(userid, true) // make jwt with user id and auth true
//         c.JSON(http.StatusOK, gin.H{ //make header with token in and send
// 				"token": token,
// 		})
//     } else {
//         fmt.Println("username or password is not a match")
//         c.IndentedJSON(http.StatusUnauthorized, nil)
//     }
   
// }

/* -------------------------------------------------------------------------------------------------------------------------*/


// func signin (c *gin.Context) {// takes id and sets up bucket and mongodb
//     userid := c.Param("userid") // get id from url request
    
//     type usrsignin struct { 
//     Username     string  `bson:"Username" json:"Username"`
//     Password  string  `bson:"Password" json:"Password"`
//     }
    
//     var reqbody usrsignin // declare new instance of struct type
    
//     if err := c.BindJSON(&reqbody); err != nil {
//         fmt.Println(err)
//         return
//     }
 
    
//     collection := client.Database("datingapp").Collection("userdata")// connect to db and collection
    
//     result := MongoFields{}

//     ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
//     // Find the document that mathes the id from the request.
// 	// call the collection's Find() method and return object into result
//     err = collection.FindOne(ctx, bson.M{"Username": userid}).Decode(&result)

//     if result.Password == reqbody.Password {
//         fmt.Println("password matches")
//         token := makejwt(userid, true) // make jwt
//         c.JSON(http.StatusOK, gin.H{ 
// 				"token": token,
// 			     })
//     } else {
//         fmt.Println("username or password is not a match")
//         c.IndentedJSON(http.StatusUnauthorized, nil)
//     }
   
// }


/* -------------------------------------------------------------------------------------------------------------------------*/

func postfile(c *gin.Context) {// post file takes file from request form, runs upload func, puts in s3, returns s3 address.

    type formData struct {
        User     string `form:"user" binding:"required"`
        Type string `form:"type" binding:"required"`
    }


    file, header, err := c.Request.FormFile("file") // get file from request body
    if err != nil {
        fmt.Println(err)
    }
 
   filename := header.Filename
    
   fileread, err := ioutil.ReadAll(file) // read the file to variable 

    userid := c.Param("userid") // get id from url request
    requestType := c.PostForm("type")

    collection := client.Database("datingapp").Collection("userdata")

    imgaddress := uploaditem(userid, filename, fileread)// call upload func returns uploaded img url
    
    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request.
    err = collection.FindOne(ctx, bson.M{"Username": userid}).Decode(&currentDoc)
    //create filter 
    filter := bson.M{"Username": userid}
    update := bson.D{}
    if requestType == "profPic" {
        currentDoc.Profpic= imgaddress //replace url to profpic section of user object
        update = bson.D{
            {"$set", bson.D{
                {"Profpic", currentDoc.Profpic},
            }},
        }
    } else if requestType == "addPhotos" {
        currentDoc.Photos= append(currentDoc.Photos, imgaddress) //append to list of users photo urls
        update = bson.D{
            {"$set", bson.D{
                {"Photos", currentDoc.Photos},
            }},
        }
    }
    
    //put to db
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    c.IndentedJSON(http.StatusOK, currentDoc)
}


/* -------------------------------------------------------------------------------------------------------------------------*/

func stdget(c *gin.Context) {// gets comments for a specified user/ all users if on home feed page
    
    type qrystruct struct { 
      Page     string  `bson:"Page" json:"Page"`
      UserName string   `bson:"UserName" json:"UserName"`
    }
    
    var qry qrystruct 
    if err := c.Bind(&qry); err != nil {
        fmt.Println(err)
        return
    }
 
    
    collection := client.Database("datingapp").Collection("userdata")// connect to db and collection.

    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

    var allPosts []PostData
    
    if qry.Page == "all" { // used when home feed wants all comments from all users
        findOptions := options.Find()
        findOptions.SetLimit(2)

        var results []*MongoFields

        // Passing bson.D{{}} as the filter matches all documents in the collection
        cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
        if err != nil {
            log.Fatal(err)
        }

        // Iterating through cursor decode documents one at a time
        for cur.Next(context.TODO()) {
            
            var elem MongoFields
            err := cur.Decode(&elem)
            if err != nil {
                log.Fatal(err)
            }

            results = append(results, &elem)
        }

        if err := cur.Err(); err != nil {
            log.Fatal(err)
        }

        cur.Close(context.TODO())

        for _, record := range results {
            for _, posts := range record.Posts{
              allPosts = append(allPosts, posts)
            }
        }

        sort.Slice(allPosts, func(i, j int) bool { //Sort posts or comments by time descending
            return allPosts[i].TimeStamp > allPosts[j].TimeStamp
        })

        var json2send MongoFields
        json2send.Posts = allPosts // send in this format so front end can use same func to access single user profile requests
        err = collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)
        json2send.Profpic = currentDoc.Profpic
        json2send.Bio = currentDoc.Bio 
        json2send.Photos = currentDoc.Photos 
        
        c.IndentedJSON(http.StatusOK, json2send)
    } else if qry.Page == "media" { // if page is users media section
        err = collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)

        c.IndentedJSON(http.StatusOK, currentDoc)

    } else{

        err = collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)

        sort.Slice(currentDoc.Posts, func(i, j int) bool {
            return currentDoc.Posts[i].TimeStamp > currentDoc.Posts[j].TimeStamp
        })

        c.IndentedJSON(http.StatusOK, currentDoc)
     }
   
}


func updatebio(c *gin.Context) {// updates user bio section


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
 
    collection := client.Database("datingapp").Collection("userdata")

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


/* -------------------------------------------------------------------------------------------------------------------------*/

func postMsg(c *gin.Context) {// std post creates doc from query, finds it and returns it


    type stdputstruct struct { // create struct with data types of incoming json object- used to find out what type of query is needed
        Username     string  `bson:"Username" json:"Username"`
        Updatetype  string  `bson:"Updatetype" json:"Updatetype"`
        IsReply  string  `bson:"IsReply" json:"IsReply"`// see if msg is a reply to post
        LikeSent Likes   `bson:"LikeSent" json:"LikeSent"`
        ReplyCmt msgCmts `bson:"ReplyCmt" json:"ReplyCmt"`
        PostIndx   int     `bson:"PostIndx" json:"PostIndx"` // used to add comments to specific post
        Key2updt string     `bson:"Key2updt" json:"Key2updt"`
        Value2updt  PostData `bson:"Value2updt" json:"Value2updt"`
    }
    
    var reqbody stdputstruct
   
    
    if err := c.BindJSON(&reqbody); err != nil {
        fmt.Println(err)
        return
    }
    
    collection := client.Database("datingapp").Collection("userdata")
    
    currentDoc := MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request. call the collection's Find() method and return object into result
    err = collection.FindOne(ctx, bson.M{"Username": reqbody.Value2updt.Username}).Decode(&currentDoc)

    //create filter 
    filter := bson.M{"Username": ""}

    // find out if msg to be added was a reply to a comment, a new post or a like.
    if reqbody.IsReply == "isResponse" {
        
        for i:=0; i<len(currentDoc.Posts); i++ {
            if currentDoc.Posts[i].PostNum == reqbody.PostIndx {
                currentDoc.Posts[i].Comments= append(currentDoc.Posts[i].Comments, reqbody.ReplyCmt) // add reply to post 
                filter = bson.M{"Username": currentDoc.Posts[i].Username} // change filter to post maker
            }
        }
    } else if reqbody.IsReply == "isCmt" {
        currentDoc.LastCommentNum += 1 
        reqbody.Value2updt.PostNum = currentDoc.LastCommentNum
        currentDoc.Posts= append(currentDoc.Posts, reqbody.Value2updt) 
        filter = bson.M{"Username": reqbody.Value2updt.Username} // set filter to username of post creator
        
    } else if reqbody.IsReply == "cmtLiked" { 
        for i:=0; i<len(currentDoc.Posts); i++ {
            if currentDoc.Posts[i].PostNum == reqbody.PostIndx {
                currentDoc.Posts[i].Likes = append(currentDoc.Posts[i].Likes, reqbody.LikeSent)  // add like to post 
                filter = bson.M{"Username": currentDoc.Posts[i].Username} // change filter to post maker
            }
        }
    } 

    update := bson.D{
        {reqbody.Updatetype, bson.D{
            {reqbody.Key2updt, currentDoc.Posts},
        }},
    }
    
    //put to db
    _, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    //update post index count

    update = bson.D{
        {reqbody.Updatetype, bson.D{
            {"LastCommentNum", currentDoc.LastCommentNum},
        }},
    }
    
    //put to db
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    result := MongoFields{}

    ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)

    var allPosts []PostData
    
    //check if request is made from profile page or news feed
    if reqbody.Username == "all" {

        findOptions := options.Find()
        findOptions.SetLimit(2)

        var results []*MongoFields

        cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
        if err != nil {
            log.Fatal(err)
        }

        //Iterate through the cursor decode documents one at a time
        for cur.Next(context.TODO()) {
            
            var elem MongoFields
            err := cur.Decode(&elem)
            if err != nil {
                log.Fatal(err)
            }

            results = append(results, &elem)
        }

        if err := cur.Err(); err != nil {
            log.Fatal(err)
        }

        cur.Close(context.TODO())

        for _, record := range results {
            for _, posts := range record.Posts{
              allPosts = append(allPosts, posts)
            }
        }

        sort.Slice(allPosts, func(i, j int) bool {
            return allPosts[i].TimeStamp > allPosts[j].TimeStamp
        })

        var json2send MongoFields
        json2send.Posts = allPosts // send in this format so front end can use same function to process other requests
        err = collection.FindOne(ctx, bson.M{"Username": reqbody.Value2updt.SessionUser}).Decode(&currentDoc)
        json2send.Profpic = currentDoc.Profpic
        json2send.Bio = currentDoc.Bio 
        json2send.Photos = currentDoc.Photos 
        c.IndentedJSON(http.StatusOK, json2send)
    } else {
        err = collection.FindOne(ctx, bson.M{"Username": reqbody.Value2updt.Username}).Decode(&result)
        sort.Slice(result.Posts, func(i, j int) bool {
            return result.Posts[i].TimeStamp > result.Posts[j].TimeStamp
        })

        c.IndentedJSON(http.StatusOK, result)
    }
   
}

/* -------------------------------------------------------------------------------------------------------------------------*/

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


/* -------------------------------------------------------------------------------------------------------------------------*/

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

/* -------------------------------------------------------------------------------------------------------------------------*/

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
/* -------------------------------------------------------------------------------------------------------------------------*/

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

/* -------------------------------------------------------------------------------------------------------------------------*/

func uploaditem (bucket string, filename string, filebytes []byte) string {// upload file to s3 with the bucket name and file adress passed to it

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

/* -------------------------------------------------------------------------------------------------------------------------*/

func GetTempLoc(filename string) string { // get the temp location of files sent in post requests to be used for aws upload without having to store data locally
    return strings.TrimRight(os.TempDir(), "/") + "/" + filename
}

/*--------------------------------------------------------------------------------------------------------------------------*/

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

/*--------------------------------------------------------------------------------------------------------------------------*/
func makejwt (userid string, isauth bool) string {
   mySigningKey := []byte("AllYourBase")//base for ecoding eg private key

   type Claims struct { // for claim settings
	Username string `json:"username"`
	jwt.StandardClaims
   }

    // Create the Claims
   // Declare the expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: userid,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)// converts to json and signs/encodes
    ss, err := token.SignedString(mySigningKey)// encode to string with chosen base string
    fmt.Printf("jwt created %v %v", ss, err)
    return ss


}

/*--------------------------------------------------------------------------------------------------------------------------*/
func chkauth (c *gin.Context) { // checks for authentication

    var jwtKey = []byte("AllYourBase")
    
    type Claims struct {
        Username string `json:"username"`
        jwt.StandardClaims
    }
    
    
    type jwtdata struct { // for user sign in  data
    Data1     string  `bson:"Data1" json:"Data1"`
    }
    
    var reqbody jwtdata 
    
    if err := c.BindJSON(&reqbody); err != nil { // unmarshall json
        fmt.Println(err)
        return
    }
    
    tknStr := reqbody.Data1           
    
    // check token is valid return username in response

    // Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.IndentedJSON(http.StatusUnauthorized, nil)
			return
		}
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}
	if !tkn.Valid {
		c.IndentedJSON(http.StatusUnauthorized, nil)
		return
	}

    type authd struct {
        AuthdUser     string  `bson:"AuthdUser" json:"AuthdUser"`
    }

    var auth authd 
    
    auth.AuthdUser = claims.Username 

	// Finally, return the welcome message to the user, along with their
	// username given in the token
    c.IndentedJSON(http.StatusOK, auth)

}

/*--------------------------------------------------------------------------------------------------------------------------*/

func exitErrorf(msg string, args ...interface{}) {// hadles aws errors
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

/*--------------------------------------------------------------------------------------------------------------------------*/

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