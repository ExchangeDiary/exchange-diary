package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/kakao"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	kakaoOAuth "golang.org/x/oauth2/kakao"
)

const defaultRedirectURL = "/api/v1/authentication/authenticated"

// AuthController ...
type AuthController interface {
	Redirect() gin.HandlerFunc
	Login() gin.HandlerFunc
	Authenticate() gin.HandlerFunc
}

type authController struct {
	client         configs.Client
	kakaoOAuthConf oauth2.Config
	memberService  service.MemberService
	tokenService   service.TokenService
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
		memberService: memberService,
		tokenService:  tokenService,
	}
}

func (ac *authController) Redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectURL := ""
		switch c.Param("auth_type") {
		case kakao.AuthType:
			redirectURL = kakaoLoginURL(&ac.kakaoOAuthConf)
		}
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}

func (ac *authController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Param("auth_type") {
		case kakao.AuthType:
			ac.kakaoLogin(c)
		}

	}
}

func (ac *authController) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCode := c.Param("authCode")
		c.JSON(http.StatusOK, gin.H{
			"authCode": authCode,
		})
	}
}

func (ac *authController) kakaoLogin(c *gin.Context) {
	code := c.Query("code")
	token, err := ac.kakaoOAuthConf.Exchange(context.TODO(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	client := ac.kakaoOAuthConf.Client(context.TODO(), token)
	kakaoClient := kakao.NewClient(ac.client.Kakao.BaseURL, client)

	kakaoUser, err := kakaoClient.GetKakaoUserInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	member, err := ac.memberService.GetByEmail(kakaoUser.Account.Email)
	// When the user try to login with registered email and different authType.
	if member != nil && member.AuthType != kakao.AuthType {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}
	member, err = ac.memberService.Create(
		kakaoUser.Account.Email,
		kakaoUser.Profile.NickName,
		kakaoUser.Profile.ProfileImage,
		kakao.AuthType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	authCode, err := ac.tokenService.IssueAuthCode(member.Email, member.AuthType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://localhost:8080%s?authcode=%s", defaultRedirectURL, authCode))
}

func kakaoLoginURL(kakaoOAuth *oauth2.Config) string {
	return kakaoOAuth.AuthCodeURL("state", oauth2.AccessTypeOnline)
}
