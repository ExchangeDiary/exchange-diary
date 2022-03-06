package middleware

import (
	"net/http"
	"strings"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	authorizationHeader = "Authorization"
	tokenBearer         = "Bearer "
)

// AuthenticationFilter ...
type AuthenticationFilter struct {
	verifier service.TokenVerifier
}

// NewAuthenticationFilter ...
func NewAuthenticationFilter(verifier service.TokenVerifier) *AuthenticationFilter {
	return &AuthenticationFilter{verifier: verifier}
}

// Authenticate ...
func (f *AuthenticationFilter) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get(authorizationHeader)
		if bearerToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header"})
			c.Abort()
			return
		}
		clientToken := strings.Replace(bearerToken, tokenBearer, "", 1)
		claims, err := f.verifier.Verify(clientToken)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		currentMember := application.CurrentMemberDTO{
			ID:    claims.ID,
			Name:  claims.Name,
			Email: claims.Email,
		}
		c.Set(application.CurrentMemberKey, currentMember)

		logger.Info("current member", zap.Uint("ID", currentMember.ID), zap.String("Name", currentMember.Name), zap.String("Email", currentMember.Email))
	}
}
