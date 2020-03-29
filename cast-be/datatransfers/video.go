package datatransfers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `json:"-" bson:"_id"`
	Hash        string             `json:"hash" bson:"hash"` // used for querying
	Type        string             `json:"type" bson:"type"` // "live" or "vod"
	Title       string             `json:"title" bson:"title"`
	Author      UserItem           `json:"author" bson:"author"`
	Description string             `json:"description" bson:"description"`
	Tags        []string           `json:"tags" bson:"tags"`
	Views       int                `json:"views" bson:"views"`
	Duration    int64              `json:"duration,omitempty" bson:"duration"` // only for VODs
	IsLive      bool               `json:"is_live,omitempty" bson:"is_live"`   // only for Live
	Resolutions int                `json:"resolutions" bson:"resolutions"`     // 0:None, 1:180p, 2:360p, 3:480p, 4:720p, 5:1080p, only for VODs
	Likes       int                `json:"likes" bson:"-"`
	Liked       bool               `json:"liked" bson:"-"`
	Comments    []Comment          `json:"comments" bson:"-"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type VideoInsert struct {
	ID          primitive.ObjectID `bson:"_id"`
	Hash        string             `bson:"hash"` // used for querying
	Type        string             `bson:"type"` // "live" or "vod"
	Title       string             `bson:"title"`
	Author      primitive.ObjectID `bson:"author"`
	Description string             `bson:"description"`
	Tags        []string           `bson:"tags"`
	Views       int                `json:"views" bson:"views"`
	Duration    int64              `bson:"duration"`    // only for VODs
	IsLive      bool               `bson:"is_live"`     // only for Live
	Resolutions int                `bson:"resolutions"` // 0:None, 1:180p, 2:360p, 3:480p, 4:720p, 5:1080p, only for VODs
	CreatedAt   time.Time          `bson:"created_at"`
}

type VideoUpload struct {
	Title       string   `form:"title"`
	Description string   `form:"description"`
	TagsMerged  string   `form:"tags"`
	Tags        []string `form:"-"`
}

type VideoEdit struct {
	Hash        string   `form:"hash"`
	Title       string   `form:"title"`
	Description string   `form:"description"`
	TagsMerged  string   `form:"tags"`
	Tags        []string `form:"-"`
}

type Comment struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	Hash      string             `json:"-" bson:"hash"`
	Content   string             `json:"content" bson:"content"`
	Author    UserItem           `json:"author" bson:"author"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type CommentInsert struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Hash      string             `json:"hash" bson:"hash"`
	Content   string             `json:"content" bson:"content"`
	Author    primitive.ObjectID `json:"author" bson:"author"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Like struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Hash      string             `json:"hash" bson:"hash"`
	Author    primitive.ObjectID `json:"author" bson:"author"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Chat struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Hash      string             `json:"hash" bson:"hash"`
	Chat      string             `json:"chat" bson:"chat"`
	Author    UserItem           `json:"author" bson:"author"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type ChatInsert struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Hash      string             `json:"hash" bson:"hash"`
	Chat      string             `json:"chat" bson:"chat"`
	Author    primitive.ObjectID `json:"author" bson:"author"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
