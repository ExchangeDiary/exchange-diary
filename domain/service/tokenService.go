package service

import (
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/golang-jwt/jwt"
)

const (
	// AuthCodeSecretKey ...
	AuthCodeSecretKey = "AUTH_CODE_SECRET_KEY"
	// AccessTokenSecretKey ...
	AccessTokenSecretKey = "ACCESS_TOKEN_SECRET_KEY"
	authCodeValid        = 5 * time.Minute    // 5 minutes
	accessTokenValid     = 7 * time.Hour      // 7 hours
	refreshTokenValid    = 7 * 24 * time.Hour //  7 days
)

// TokenService ...
type TokenService interface {
	IssueAuthCode(email string, authType string) (string, error)
	IssueAccessToken(authCode string) (string, error)
	IssueRefreshToken(authCode string) (string, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type tokenService struct {
	memberService        MemberService
	authCodeVerifier     TokenVerifier
	refreshTokenVerifier TokenVerifier
}

// NewTokenService ...
func NewTokenService(service MemberService, authCodeVerifier TokenVerifier, refreshTokenVerifier TokenVerifier) TokenService {
	return &tokenService{
		memberService:        service,
		authCodeVerifier:     authCodeVerifier,
		refreshTokenVerifier: refreshTokenVerifier,
	}
}

func (s *tokenService) IssueAuthCode(email string, authType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.AuthCodeClaims{
		AuthType: authType,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(authCodeValid).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exchange-diary",
		},
	})
	return token.SignedString([]byte(s.authCodeVerifier.SecretKey))
}

func (s *tokenService) IssueAccessToken(authCode string) (string, error) {
	claims, err := s.authCodeVerifier.Verify(authCode)
	if err != nil {
		return "", err
	}
	member, err := s.memberService.GetByEmail(claims.Email)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.AuthCodeClaims{
		AuthType: member.AuthType,
		ID:       member.ID,
		Email:    member.Email,
		Name:     member.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenValid).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exchange-diary",
		},
	})
	return token.SignedString([]byte(s.authCodeVerifier.SecretKey))
}

func (s *tokenService) IssueRefreshToken(authCode string) (string, error) {
	claims, err := s.authCodeVerifier.Verify(authCode)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.AuthCodeClaims{
		AuthType: claims.AuthType,
		ID:       claims.ID,
		Email:    claims.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenValid).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exchange-diary",
		},
	})
	return token.SignedString([]byte(s.authCodeVerifier.SecretKey))
}

func (s tokenService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := s.refreshTokenVerifier.Verify(refreshToken)
	if err != nil {
		return "", err
	}
	member, err := s.memberService.GetByEmail(claims.Email)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.AuthCodeClaims{
		AuthType: member.AuthType,
		ID:       member.ID,
		Email:    member.Email,
		Name:     member.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(authCodeValid).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "exchange-diary",
		},
	})
	return token.SignedString([]byte(s.refreshTokenVerifier.SecretKey))
}
