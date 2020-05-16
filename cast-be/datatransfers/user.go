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
	Subscribers int                `json:"subscribers" bson:"subscribers"`
	Verified    bool               `json:"verified" bson:"verified"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at"`
}

type UserDetail struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Verified    bool   `json:"verified" bson:"verified"`
	Subscribers int    `json:"subscribers" bson:"subscribers"`
	Views       int    `json:"views" bson:"views"`
	Uploads     int    `json:"uploads" bson:"uploads"`
}

type UserItem struct {
	ID          primitive.ObjectID `json:"-" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Username    string             `json:"username" bson:"username"`
	Subscribers int                `json:"subscribers" bson:"subscribers"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserVerify struct {
	Key string `json:"key"`
}

type UserResend struct {
	Email string `json:"email"`
}

type UserFieldCheck struct {
	Field string
	Value string
}

type UserEditForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UserUpdatePassword struct {
	UserVerify
	Password string `json:"password"`
}
