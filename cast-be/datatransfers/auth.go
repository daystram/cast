package datatransfers

// JWT claims struct
type JWTClaims struct {
	ID       string
	Expiry   int64
	Remember bool
}
