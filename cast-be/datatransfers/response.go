package datatransfers

// API response wrapper
type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
	Code  int         `json:"code"`
}
