package constants

import "time"

const (
	ContextParamUserID      = "user_id"
	AuthenticationTimeout   = 48 * time.Hour
	AuthenticationCookieKey = "user"
)
