package datatransfers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"-" bson:"password"`
	Subscribers int                `json:"-" bson:"subscribers"`
	//Subscribing []primitive.ObjectID // TODO: normalize to other table
	//Likes       []primitive.ObjectID // TODO: normalize to other table
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
}

type UserDetail struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Subscribers int    `json:"subscribers" bson:"subscribers"`
	Views       int    `json:"views" bson:"views"`
	Uploads     int    `json:"uploads" bson:"uploads"`
}

type UserItem struct {
	Name        string `json:"name" bson:"name"`
	Username    string `json:"username" bson:"username"`
	Subscribers int    `json:"subscribers" bson:"subscribers"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserFieldCheck struct {
	Field string
	Value string
}

type UserEdit struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
