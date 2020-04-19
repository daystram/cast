package datatransfers

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}
