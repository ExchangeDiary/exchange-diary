package entity

import "github.com/golang-jwt/jwt"

type AuthCodeClaims struct {
	AuthType string `json:"auth_type"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
