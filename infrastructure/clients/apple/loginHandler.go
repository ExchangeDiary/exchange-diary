package apple

import (
	"context"
	"fmt"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	appleOAuth "github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type LoginHandler struct {
	Conf          configs.Apple
	MemberService service.MemberService
	TokenService  service.TokenService
}

func (a LoginHandler) Handle(authorizationCode string, identityToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey, err := ioutil.ReadFile(a.Conf.Oauth.KeyPath)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// validate identity token
		err = verifyIdToken(a.Conf.Oauth.ClientID, identityToken)
		if err != nil {
			c.String(http.StatusUnauthorized, "identity token verification failed")
			c.Abort()
		}

		secret, err := appleOAuth.GenerateClientSecret(string(authKey), a.Conf.Oauth.TeamId, a.Conf.Oauth.ClientID, a.Conf.Oauth.KeyID)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		client := appleOAuth.New()

		req := appleOAuth.AppValidationTokenRequest{
			ClientID:     a.Conf.Oauth.ClientID,
			ClientSecret: secret,
			Code:         authorizationCode,
		}
		var resp appleOAuth.ValidationResponse

		err = client.VerifyAppToken(context.Background(), req, &resp)
		if err != nil {
			c.String(http.StatusUnauthorized, "verification failed")
			c.Abort()
			return
		}

		claim, _ := appleOAuth.GetClaims(resp.IDToken)

		email := fmt.Sprintf("%v", (*claim)["email"])

		member, err := a.MemberService.GetByEmail(email)

		// When the user try to login with registered email and different authType.
		if member != nil && member.AuthType != AuthType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The account is already registered"})
			return
		}

		authCode, err := a.TokenService.IssueAuthCode(email, AuthType)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"auth_code": authCode})
	}
}
