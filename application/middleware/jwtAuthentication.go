package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set(application.CurrentMemberKey, application.CurrentMemberDTO{
			ID:    claims.ID,
			Name:  claims.Name,
			Email: claims.Email,
		})
		fmt.Printf("[Current Member]: %+v", c.MustGet(application.CurrentMemberKey))
	}
}
