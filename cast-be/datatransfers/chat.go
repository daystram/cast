package datatransfers

import "time"

type ChatOutgoing struct {
	Author    string    `json:"author"`
	Chat      string    `json:"chat"`
	CreatedAt time.Time `json:"created_at"`
}
