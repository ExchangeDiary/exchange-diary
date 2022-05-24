package google

import (
	"context"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	googleApiOauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"net/http"
)

type LoginHandler struct {
	OauthConf     oauth2.Config
	MemberService service.MemberService
	TokenService  service.TokenService
}

func (g LoginHandler) Handle(authorizationCode string, identityToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := g.OauthConf.Exchange(context.Background(), authorizationCode)

		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		client := g.OauthConf.Client(context.Background(), token)
		googleService, err := googleApiOauth.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		googleUser, err := googleService.Userinfo.Get().Do()
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		member, err := g.MemberService.GetByEmail(googleUser.Email)

		// When the user try to login with registered email and different authType.
		if member != nil && member.AuthType != AuthType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The account is already registered"})
			return
		}

		authCode, err := g.TokenService.IssueAuthCode(googleUser.Email, AuthType)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"auth_code": authCode})
	}
}
