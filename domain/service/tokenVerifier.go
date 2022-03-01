package service

import (
	"fmt"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/golang-jwt/jwt"
)

// TokenVerifier ...
type TokenVerifier struct {
	SecretKey string
}

// NewTokenVerifier ...
func NewTokenVerifier(secretKey string) TokenVerifier {
	return TokenVerifier{SecretKey: secretKey}
}

// Verify ...
func (t *TokenVerifier) Verify(authCode string) (claims *entity.AuthCodeClaims, err error) {
	token, err := jwt.ParseWithClaims(
		authCode,
		&entity.AuthCodeClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(t.SecretKey), nil
		},
	)
	claims, ok := token.Claims.(*entity.AuthCodeClaims)
	if !ok {
		err = fmt.Errorf("the token is invalid")
		return
	}
	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		fmt.Println(claims.ExpiresAt)
		err = fmt.Errorf("token is expired")
		return
	}
	
	return claims, err
}
