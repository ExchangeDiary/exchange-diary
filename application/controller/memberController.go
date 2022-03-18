package controller

import (
	"net/http"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// MemberController ...
type MemberController interface {
	Get() gin.HandlerFunc
	Post() gin.HandlerFunc
	Patch() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

type memberController struct {
	memberService service.MemberService
}

type memberRequest struct {
	Email      string `json:"email"`
	Name       string `json:"name,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
	AuthType   string `json:"auth_type,omitempty"`
	AlarmFlag  bool   `json:"alarm_flag,omitempty"`
}

type memberResponse struct {
	Email      string `json:"email"`
	Name       string `json:"name,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
	AuthType   string `json:"auth_type"`
	AlarmFlag  bool   `json:"alarm_flag,omitempty"`
}

// NewMemberController ...
func NewMemberController(memberService service.MemberService) MemberController {
	return &memberController{memberService: memberService}
}

// @Summary Member 조회
// @Description	 email 주소를 통해 가입된 member를 조회한다.
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        email   path   string  true "사용자 이메일"
// @Success      200     {object}  memberResponse
// @Failure      400
// @Failure      500
// @Router       /member/{email} [get]
// @Security ApiKeyAuth
func (mc *memberController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		member, err := mc.memberService.GetByEmail(email)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := memberResponse{
			Email:      member.Email,
			Name:       member.Name,
			ProfileURL: member.ProfileURL,
			AuthType:   member.AuthType,
			AlarmFlag:  member.AlarmFlag,
		}
		c.JSON(http.StatusOK, response)
	}
}

// @Summary Member 생성
// @Description	 member를 새로 생성한다.
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        member   body  memberRequest  true "member 생성 요청 body"
// @Success      201      {object}  memberResponse
// @Failure      400
// @Failure      500
// @Router       /member [post]
// @Security ApiKeyAuth
func (mc *memberController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request memberRequest
		if err := c.BindJSON(&request); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		member, err := mc.memberService.Create(request.Email, request.Name, request.ProfileURL, request.AuthType)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		response := memberResponse{
			Email:      member.Email,
			Name:       member.Name,
			ProfileURL: member.ProfileURL,
			AuthType:   member.AuthType,
			AlarmFlag:  member.AlarmFlag,
		}
		c.JSON(http.StatusCreated, response)
	}
}

// @Summary Member 수정
// @Description	 해당 member를 수정한다.
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        member   body  memberRequest  true "member 수정 요청 body"
// @Success      200      {object}  memberResponse
// @Failure      400
// @Failure      500
// @Router       /member [patch]
// @Security ApiKeyAuth
func (mc *memberController) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request memberRequest
		if err := c.BindJSON(&request); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		member, err := mc.memberService.GetByEmail(request.Email)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		patchedMember, err := mc.memberService.Update(patchMember(member, request))
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		response := memberResponse{
			Email:      patchedMember.Email,
			Name:       patchedMember.Name,
			ProfileURL: patchedMember.ProfileURL,
			AuthType:   patchedMember.AuthType,
			AlarmFlag:  patchedMember.AlarmFlag,
		}
		c.JSON(http.StatusOK, response)
	}
}

func patchMember(original *entity.Member, patch memberRequest) *entity.Member {
	isIdentical := true
	if patch.Name != "" {
		original.Name = patch.Name
		isIdentical = false
	}
	if patch.ProfileURL != "" {
		original.ProfileURL = patch.ProfileURL
		isIdentical = false
	}
	if patch.AlarmFlag != original.AlarmFlag {
		original.AlarmFlag = patch.AlarmFlag
		isIdentical = false
	}
	if isIdentical == false {
		original.UpdatedAt = time.Now()
	}
	return original
}

// @Summary Member 삭제
// @Description	 해당 member를 삭제한다.
// @Tags         members
// @Accept       json
// @Produce      json
// @Param        email   path  string  true "사용자 이메일"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /member/{email} [delete]
// @Security ApiKeyAuth
func (mc *memberController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		err := mc.memberService.Delete(email)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.Status(http.StatusOK)
	}
}
