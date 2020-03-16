package datatransfers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Video struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Type        string             // "live" or "vod"
	Title       string
	Author      User
	Description string
	Duration    int64  `json:"omitempty"` // only for VODs
	IsLive      bool   `json:"omitempty"` // only for Live
	Resolutions string `json:"omitempty"` // 0:None, 1:180p, 2:360p, 3:480p, 4:720p, 5:1080p, only for VODs
	CreatedAt   time.Time
}
