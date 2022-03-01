package controller

import (
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/gin-gonic/gin"
)

// TokenController ...
type TokenController interface {
	GetToken() gin.HandlerFunc
	RefreshAccessToken() gin.HandlerFunc
}

type tokenController struct {
	service service.TokenService
}

// TokenRequest ...
type TokenRequest struct {
	AuthCode string `json:"authCode"`
}

// TokenRefreshRequest ...
type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// NewTokenController ...
func NewTokenController(service service.TokenService) TokenController {
	return &tokenController{service: service}
}

func (tc *tokenController) GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody TokenRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accessToken, err := tc.service.IssueAccessToken(requestBody.AuthCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errror": err.Error()})
			return
		}
		refreshToken, err := tc.service.IssueRefreshToken(requestBody.AuthCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token := entity.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		c.JSON(http.StatusOK, token)
	}
}

func (tc tokenController) RefreshAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody TokenRefreshRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accessToken, err := tc.service.RefreshAccessToken(requestBody.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errror": err.Error()})
			return
		}
		token := entity.Token{
			AccessToken:  accessToken,
			RefreshToken: requestBody.RefreshToken,
		}
		c.JSON(http.StatusOK, token)
	}
}
