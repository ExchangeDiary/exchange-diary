package controller

import (
	"context"
	"fmt"
	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/kakao"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	googleOAuth "golang.org/x/oauth2/google"
	kakaoOAuth "golang.org/x/oauth2/kakao"
	googleApiOauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"net/http"
)

const defaultRedirectURL = "/v1/authentication/authenticated"

// AuthController ...
type AuthController interface {
	Redirect() gin.HandlerFunc
	Login() gin.HandlerFunc
	Authenticate() gin.HandlerFunc
	MockRegister() gin.HandlerFunc
}

type authController struct {
	client          configs.Client
	kakaoOAuthConf  oauth2.Config
	googleOAuthConf oauth2.Config
	memberService   service.MemberService
	tokenService    service.TokenService
}

// NewAuthController ...
func NewAuthController(client configs.Client, memberService service.MemberService, tokenService service.TokenService) AuthController {
	return &authController{
		client: client,
		kakaoOAuthConf: oauth2.Config{
			ClientID:     client.Kakao.Oauth.ClientID,
			ClientSecret: client.Kakao.Oauth.ClientSecret,
			Endpoint:     kakaoOAuth.Endpoint,
			RedirectURL:  client.Kakao.Oauth.RedirectURL,
		},
		googleOAuthConf: oauth2.Config{
			ClientID:     client.Google.Oauth.ClientID,
			ClientSecret: client.Google.Oauth.ClientSecret,
			Endpoint:     googleOAuth.Endpoint,
			RedirectURL:  client.Google.Oauth.RedirectURL,
			Scopes: []string{
				googleApiOauth.UserinfoEmailScope,
			},
		},
		memberService: memberService,
		tokenService:  tokenService,
	}
}

// @Summary      LoginForm Redirect
// @Description	 요청한 auth_type의 login form URL로 redirect됩니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth_type  path string  true  "kakao | google | apple"
// @Success      301
// @Failure      400
// @Failure      500
// @Router       /login/{auth_type} [get]
func (ac *authController) Redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectURL := ""
		switch c.Param("auth_type") {
		case kakao.AuthType:
			redirectURL = kakaoLoginURL(&ac.kakaoOAuthConf)
		case google.AuthType:
			redirectURL = googleLoginURL(&ac.googleOAuthConf)
		}
		logger.Info("Redirect login is accepted", zap.String("redirectURL", redirectURL))
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}

// @Summary      Login
// @Description	 해당 auth_type의 login form을 기입 후, redirect되는 곳으로 기입한 email 그리고 요청한 auth_type에 대한 auth code를 발급한다.
// @Description  클라이언트 URL에 query string으로 auth code가 담긴 채로 redirect하여 로그인 및 회원가입이 진행되도록 한다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth_type   path      string  true  "kakao | google | apple"
// @Success      301
// @Failure      400
// @Failure      500
// @Router       /authentication/login/{auth_type} [get]
func (ac *authController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Param("auth_type") {
		case kakao.AuthType:
			ac.kakaoLogin(c)
		case google.AuthType:
			ac.googleLogin(c)
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

type mockMemberRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type mockMemberResponse struct {
	AuthCode     string `json:"authCode"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary      (debug) mock login / register
// @Description	 클라 테스트용. 주어진 email이 db에 없으면 회원가입 프로세스 동시에 진행
// @Description	 AccessToken을 사용해서 헤더에 {"Authorization": AccessToken} 넣어주면 된다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        member  body     mockMemberRequest  true  "모킹할 유저정보"
// @Success      200  {object}   mockMemberResponse
// @Failure      400
// @Failure      500
// @Router       /authentication/mock [post]
func (ac *authController) MockRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req mockMemberRequest
		if err := c.BindJSON(&req); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		member, err := ac.memberService.GetByEmail(req.Email)
		if err != nil {
			defaultProfileURL := "https://user-images.githubusercontent.com/37536298/153554715-f821d0f8-8f51-4f4c-b9e6-a19e02ecb5c2.png"
			member, err = ac.memberService.Create(
				req.Email,
				req.Name,
				defaultProfileURL,
				kakao.AuthType,
			)
			if err != nil {
				logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		authCode, err := ac.tokenService.IssueAuthCode(member.Email, member.AuthType)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		accessToken, err := ac.tokenService.IssueAccessToken(authCode)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		refreshToken, err := ac.tokenService.IssueRefreshToken(authCode)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, mockMemberResponse{
			AuthCode:     authCode,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (ac *authController) kakaoLogin(c *gin.Context) {
	code := c.Query("code")
	logger.Info("kakao login code: " + code)
	token, err := ac.kakaoOAuthConf.Exchange(context.Background(), code)

	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	client := ac.kakaoOAuthConf.Client(context.Background(), token)
	kakaoClient := kakao.NewClient(ac.client.Kakao.BaseURL, client)

	kakaoUser, err := kakaoClient.GetKakaoUserInfo()
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	member, err := ac.memberService.GetByEmail(kakaoUser.Account.Email)

	// When the user try to login with registered email and different authType.
	if member != nil && member.AuthType != kakao.AuthType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The account is already registered"})
		return
	}

	authCode, err := ac.tokenService.IssueAuthCode(kakaoUser.Account.Email, kakao.AuthType)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: domain 주소 env로 빼기
	//c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://localhost:8080%s?%s=%s", defaultRedirectURL, application.AuthCodeKey, authCode))
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("https://exchange-diary-b4mzhzbzcq-du.a.run.app%s?%s=%s", defaultRedirectURL, application.AuthCodeKey, authCode))
}

func (ac *authController) googleLogin(c *gin.Context) {
	code := c.Query("code")
	logger.Info(fmt.Sprintf("google login code is < %s >", code))
	token, err := ac.googleOAuthConf.Exchange(context.Background(), code)

	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	client := ac.googleOAuthConf.Client(context.Background(), token)
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

	member, err := ac.memberService.GetByEmail(googleUser.Email)

	// When the user try to login with registered email and different authType.
	if member != nil && member.AuthType != google.AuthType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The account is already registered"})
		return
	}

	authCode, err := ac.tokenService.IssueAuthCode(googleUser.Email, google.AuthType)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: domain 주소 env로 빼기
	//c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://localhost:8080%s?%s=%s", defaultRedirectURL, application.AuthCodeKey, authCode))
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("https://exchange-diary-b4mzhzbzcq-du.a.run.app%s?%s=%s", defaultRedirectURL, application.AuthCodeKey, authCode))
}

func kakaoLoginURL(kakaoOAuth *oauth2.Config) string {
	return kakaoOAuth.AuthCodeURL("state", oauth2.AccessTypeOnline)
}

func googleLoginURL(googleOauth *oauth2.Config) string {
	return googleOauth.AuthCodeURL("state", oauth2.AccessTypeOnline)
}
