package constants

import "time"

const (
	ContextParamUserID            = "user_id"
	AuthenticationTimeout         = 24 * time.Hour
	AuthenticationTimeoutExtended = 365 * 24 * time.Hour
	AuthenticationCookieKey       = "cast"
)
