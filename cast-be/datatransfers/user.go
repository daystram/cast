package datatransfers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id"`
	Name        string
	Username    string
	Email       string
	//Subscribing []primitive.ObjectID // TODO: normalize to other table
	//Likes       []primitive.ObjectID // TODO: normalize to other table
	CreatedAt   time.Time
}

type UserLogin struct {
	UsernameEmail string `json:"user"`
	Password      string `json:"password"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
