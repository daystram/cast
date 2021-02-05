package datatransfers

import (
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
}

type UserDetail struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Subscribers int    `json:"subscribers" bson:"subscribers"`
	Views       int    `json:"views" bson:"views"`
	Uploads     int    `json:"uploads" bson:"uploads"`
}

type UserItem struct {
	ID          string `json:"-" bson:"_id"`
	Username    string `json:"username" bson:"username"`
	Name        string `json:"name" bson:"name"`
	Subscribers int    `json:"subscribers" bson:"subscribers"`
}

type UserRegister struct {
	IDToken string `json:"id_token"`
}
