package datatransfers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Hash        string             `json:"hash" bson:"hash"` // used for querying
	Type        string             `json:"type" bson:"type"` // "live" or "vod"
	Title       string             `json:"title" bson:"title"`
	Author      UserItem           `json:"author"`
	Description string             `json:"description" bson:"description"`
	Views       int                `json:"views" bson:"views"`
	Duration    int64              `json:"duration,omitempty" bson:"duration"`       // only for VODs
	IsLive      bool               `json:"is_live,omitempty" bson:"is_live"`         // only for Live
	Resolutions int                `json:"resolutions" bson:"resolutions"` // 0:None, 1:180p, 2:360p, 3:480p, 4:720p, 5:1080p, only for VODs
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type VideoInsert struct {
	ID          primitive.ObjectID `bson:"_id"`
	Hash        string             `bson:"hash"` // used for querying
	Type        string             `bson:"type"` // "live" or "vod"
	Title       string             `bson:"title"`
	Author      primitive.ObjectID `bson:"author"`
	Description string             `bson:"description"`
	Views       int                `json:"views" bson:"views"`
	Duration    int64              `bson:"duration"`    // only for VODs
	IsLive      bool               `bson:"is_live"`     // only for Live
	Resolutions int                `bson:"resolutions"` // 0:None, 1:180p, 2:360p, 3:480p, 4:720p, 5:1080p, only for VODs
	CreatedAt   time.Time          `bson:"created_at"`
}

type VideoUpload struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Tags        string `form:"tags"`
}
