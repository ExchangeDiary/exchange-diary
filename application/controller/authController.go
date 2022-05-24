package controller

import (
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/apple"
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/kakao"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	googleOAuth "golang.org/x/oauth2/google"
	googleApiOauth "google.golang.org/api/oauth2/v2"
)

type AuthController interface {
	Login() gin.HandlerFunc
	Authenticate() gin.HandlerFunc
}

type LoginHandler interface {
	Handle(authorizationCode string, identityToken string) gin.HandlerFunc
}

type authController struct {
	kakaoLoginHandler  kakao.LoginHandler
	googleLoginHandler google.LoginHandler
	appleLoginHandler  apple.LoginHandler
}

type loginRequest struct {
	AuthorizationCode string `json:"authorizationCode"`
	IdentityToken     string `json:"identityToken,omitempty"`
}

// @Summary      Login
// @Description	 해당 auth_type과 각 vendor의 token이 주어지면 email과 요청한 auth_type에 대한 auth code를 발급한다.
// @Description  kako, google: authCode
// @Description  apple: authCode, identityToken
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth_type   path      string  true  "kakao | google | apple"
// @Success      301
// @Failure      400
// @Failure      500
// @Router       /authentication/login/{auth_type} [post]
func (ac *authController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.BindJSON(&req); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.AuthorizationCode == "" {
			c.String(http.StatusForbidden, "No authorization code provided")
			c.Abort()
			return
		}

		switch c.Param("auth_type") {
		case kakao.AuthType:
			ac.kakaoLoginHandler.Handle(req.AuthorizationCode, "")
			return
		case google.AuthType:
			ac.googleLoginHandler.Handle(req.AuthorizationCode, "")
			return
		case apple.AuthType:
			if req.IdentityToken == "" {
				c.String(http.StatusForbidden, "No identity token provided")
				c.Abort()
				return
			}
			ac.appleLoginHandler.Handle(req.AuthorizationCode, req.IdentityToken)
			return
		}
	}
}

func (ac *authController) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCode := c.Query(application.AuthCodeKey)
		c.JSON(http.StatusOK, gin.H{
			"authCode": authCode,
		})
	}
}

func NewAuthController(client configs.Client, memberService service.MemberService, tokenService service.TokenService) AuthController {
	return &authController{
		kakaoLoginHandler: kakao.LoginHandler{
			Conf:          client.Kakao,
			MemberService: memberService,
			TokenService:  tokenService,
		},
		googleLoginHandler: google.LoginHandler{
			OauthConf: oauth2.Config{
				ClientID:     client.Google.Oauth.ClientID,
				ClientSecret: client.Google.Oauth.ClientSecret,
				Endpoint:     googleOAuth.Endpoint,
				RedirectURL:  client.Google.Oauth.RedirectURL,
				Scopes: []string{
					googleApiOauth.UserinfoEmailScope,
				},
			},
			MemberService: memberService,
			TokenService:  tokenService,
		},
		appleLoginHandler: apple.LoginHandler{
			Conf:          client.Apple,
			MemberService: memberService,
			TokenService:  tokenService,
		},
	}
}
