package datatransfers

import "time"

// WS chat message
type ChatOutgoing struct {
	Author    string    `json:"author"`
	Chat      string    `json:"chat"`
	CreatedAt time.Time `json:"created_at"`
}
