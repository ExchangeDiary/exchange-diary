package middleware

import (
	"github.com/exchange-diary/domain/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AUTHORIZATION_HEADER = "Authorization"
	TOKEN_BEARER         = "Bearer "
)

type AuthenticationFilter struct {
	verifier service.TokenVerifier
}

func NewAuthenticationFilter(verifier service.TokenVerifier) *AuthenticationFilter {
	return &AuthenticationFilter{verifier: verifier}
}

func (f *AuthenticationFilter) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get(AUTHORIZATION_HEADER)
		if bearerToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header"})
			c.Abort()
			return
		}
		clientToken := strings.Replace(bearerToken, TOKEN_BEARER, "", 1)
		claims, err := f.verifier.Verify(clientToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
	}
}
