package controller

import (
	"net/http"
	"time"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// AlarmController ...
type AlarmController interface {
	List() gin.HandlerFunc
}

type alarmController struct {
	alarmService  service.AlarmService
	memberService service.MemberService
}

// NewAlarmController ...
func NewAlarmController(as service.AlarmService, ms service.MemberService) AlarmController {
	return &alarmController{
		alarmService:  as,
		memberService: ms,
	}
}

type responseAlarm struct {
	Code     string     `json:"code"`
	Title    string     `json:"title"`
	RoomName string     `json:"roomName"`
	AlarmAt  *time.Time `json:"alarmAt"`
	Author   string     `json:"author"`
}

// TODO: order by datetime
type AlarmListResponse struct {
	Alarms []responseAlarm `json:"alarms"`
}

// @Summary      List alarms
// @Description  현재 로그인한 사용자의 알람 리스트
// @Tags         alarms
// @Accept       json
// @Produce      json
// @Success      200  {object}   AlarmListResponse
// @Failure      400
// @Failure      401
// @Router       /alarms [get]
// @Security ApiKeyAuth
func (ac *alarmController) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		alarms, err := ac.alarmService.GetAll(currentMember.ID)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rep := []responseAlarm{}
		for _, a := range *alarms {
			rep = append(rep, responseAlarm{
				Code:     a.Code,
				Title:    a.Title,
				RoomName: a.RoomName,
				AlarmAt:  a.CreatedAt,
				Author:   a.Author,
			})
		}
		c.JSON(http.StatusOK, AlarmListResponse{Alarms: rep})
	}
}
