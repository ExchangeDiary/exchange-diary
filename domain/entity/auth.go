package entity

import "github.com/golang-jwt/jwt"

// AuthCodeClaims ...
type AuthCodeClaims struct {
	AuthType string `json:"auth_type"`
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

// Token ...
type Token struct {
	AccessToken  string
	RefreshToken string
}
