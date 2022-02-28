package service

import (
	"fmt"
	"github.com/exchange-diary/domain/entity"
	"github.com/golang-jwt/jwt"
)

type TokenVerifier struct {
	SecretKey string
}

func NewTokenVerifier(secretKey string) TokenVerifier {
	return TokenVerifier{SecretKey: secretKey}
}

func (t *TokenVerifier) Verify(authCode string) (*entity.AuthCodeClaims, error) {
	token, err := jwt.Parse(authCode, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return t.SecretKey, nil
	})
	claims, ok := token.Claims.(entity.AuthCodeClaims)

	if ok && token.Valid {
		return &claims, nil
	}
	return nil, err
}
