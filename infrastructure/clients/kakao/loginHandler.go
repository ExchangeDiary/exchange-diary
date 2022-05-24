package kakao

import (
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginHandler struct {
	Conf          configs.Kakao
	MemberService service.MemberService
	TokenService  service.TokenService
}

func (k LoginHandler) Handle(authorizationCode string, identityToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		kakaoClient := NewClient(k.Conf.BaseURL, authorizationCode)

		kakaoUser, err := kakaoClient.GetKakaoUserInfo()
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		member, err := k.MemberService.GetByEmail(kakaoUser.Account.Email)

		// When the user try to login with registered email and different authType.
		if member != nil && member.AuthType != AuthType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The account is already registered"})
			return
		}

		authCode, err := k.TokenService.IssueAuthCode(kakaoUser.Account.Email, AuthType)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"auth_code": authCode})
	}
}
