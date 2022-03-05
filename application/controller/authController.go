package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/kakao"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	kakaoOAuth "golang.org/x/oauth2/kakao"
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

// @Summary      login
// @Description	 회원가입 하지 않았을 경우, email로 회원가입 자동 진행
// @Description	 이후 jwt 토큰 발급에 필요한 authCode를 전달한다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth_type    query     string  true  "kakao | google | apple"  Format(string)
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /authentication/login/:auth_type [get]
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
// @Param        room  body     mockMemberRequest  true  "모킹할 유저정보"
// @Success      200  {object}   mockMemberResponse
// @Failure      400
// @Failure      500
// @Router       /authentication/mock [post]
func (ac *authController) MockRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req mockMemberRequest
		if err := c.BindJSON(&req); err != nil {
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		authCode, err := ac.tokenService.IssueAuthCode(member.Email, member.AuthType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		accessToken, err := ac.tokenService.IssueAccessToken(authCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		refreshToken, err := ac.tokenService.IssueRefreshToken(authCode)
		if err != nil {
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
	token, err := ac.kakaoOAuthConf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	client := ac.kakaoOAuthConf.Client(context.Background(), token)
	kakaoClient := kakao.NewClient(ac.client.Kakao.BaseURL, client)

	kakaoUser, err := kakaoClient.GetKakaoUserInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	member, err := ac.memberService.GetByEmail(kakaoUser.Account.Email)
	if err != nil {
		member, err = ac.memberService.Create(
			kakaoUser.Account.Email,
			kakaoUser.Profile.NickName,
			kakaoUser.Profile.ProfileImage,
			kakao.AuthType,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// When the user try to login with registered email and different authType.
	if member.AuthType != kakao.AuthType {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}

	authCode, err := ac.tokenService.IssueAuthCode(member.Email, member.AuthType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://localhost:8080%s?%s=%s", defaultRedirectURL, application.AuthCodeKey, authCode))
}

func kakaoLoginURL(kakaoOAuth *oauth2.Config) string {
	return kakaoOAuth.AuthCodeURL("state", oauth2.AccessTypeOnline)
}
