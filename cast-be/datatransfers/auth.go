package datatransfers

type JWTClaims struct {
	ID       string
	Expiry   int64
	Remember bool
}
