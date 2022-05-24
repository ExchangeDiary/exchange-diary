package apple

// AuthType value for apple
const AuthType = "apple"

type JWTTokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type JWTTokenBody struct {
	Iss            string `json:"iss"`
	Iat            int64  `json:"iat"`
	Exp            int64  `json:"exp"`
	Aud            string `json:"aud"`
	Sub            string `json:"sub"`
	AtHash         string `json:"at_hash"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	IsPrivateEmail string `json:"is_private_email"`
	RealUserStatus int64  `json:"real_user_status"`
	AuthTime       int64  `json:"auth_time"`
	Nonce          string `json:"nonce"`
}
