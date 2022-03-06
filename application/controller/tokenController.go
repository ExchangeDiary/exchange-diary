package controller

import (
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
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

type tokenResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary      JWT 토큰 발급
// @Description	 AuthCode를 전달하여, access & refresh 토큰을 발급 받는다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        room  body     TokenRequest  true  "발급받은 AuthCode"
// @Success      200  {object}   tokenResponse
// @Failure      400
// @Failure      500
// @Router       /token [post]
func (tc *tokenController) GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody TokenRequest
		if err := c.BindJSON(&requestBody); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accessToken, err := tc.service.IssueAccessToken(requestBody.AuthCode)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		refreshToken, err := tc.service.IssueRefreshToken(requestBody.AuthCode)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
	}
}

// @Summary      AccessToken 재발급
// @Description	 refresh token을 전달하여, accessToken을 재발급받는다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        room  body     TokenRefreshRequest  true  "refresh 토큰"
// @Success      200  {object}   tokenResponse
// @Failure      400
// @Failure      500
// @Router       /token/refresh [get]
func (tc tokenController) RefreshAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody TokenRefreshRequest
		if err := c.BindJSON(&requestBody); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accessToken, err := tc.service.RefreshAccessToken(requestBody.RefreshToken)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tokenResponse{AccessToken: accessToken, RefreshToken: requestBody.RefreshToken})
	}
}
