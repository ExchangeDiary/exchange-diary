package controller

import (
	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenController interface {
	GetToken() gin.HandlerFunc
	RefreshAccessToken() gin.HandlerFunc
}

type tokenController struct {
	service service.TokenService
}

type TokenRequest struct {
	authCode string `json:"auth_code"`
}

type TokenRefreshRequest struct {
	refreshToken string `json:"refresh_token"`
}

func NewTokenController(service service.TokenService) TokenController {
	return &tokenController{service: service}
}

func (tc *tokenController) GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody TokenRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		accessToken, err := tc.service.IssueAccessToken(requestBody.authCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errror": err})
		}
		refreshToken, err := tc.service.IssueRefreshToken(requestBody.authCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		accessToken, err := tc.service.RefreshAccessToken(requestBody.refreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errror": err})
		}
		token := entity.Token{
			AccessToken:  accessToken,
			RefreshToken: requestBody.refreshToken,
		}
		c.JSON(http.StatusOK, token)
	}
}
