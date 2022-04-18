package social

	type MsgCmts struct {
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
	Comments [] MsgCmts  `bson:"Comments" json:"Comments"`
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