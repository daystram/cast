package datatransfers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Invoker   primitive.ObjectID `json:"invoker" bson:"invoker"`
	Purpose   string             `json:"purpose" bson:"purpose"`
	Hash      string             `json:"hash" bson:"hash"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
}
