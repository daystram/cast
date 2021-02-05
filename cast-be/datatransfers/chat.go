package datatransfers

import "time"

// WS chat message
type ChatOutgoing struct {
	Author    string    `json:"author"`
	Chat      string    `json:"chat"`
	CreatedAt time.Time `json:"created_at"`
}

// WS notification message
type NotificationOutgoing struct {
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}
