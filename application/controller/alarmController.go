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

// AlarmListResponse ...
type AlarmListResponse struct {
	Alarms []responseAlarm `json:"alarms"`
}

// @Summary      List alarms
// @Description  현재 로그인한 사용자의 알람 리스트
// @Description  * 교환일기 방의 주기에 따라 알람은 삭제된다. 즉 턴이 변경되면 알림들 제거 된다.
// @Description  * 실제 구현은 알람 생성 시 (room & member * code)유니크 검사를 통해서 중복 되는 알림코드는 제거 함으로써, 알람 row 숫자를 관리한다.
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
				AlarmAt:  a.AlarmAt,
				Author:   a.Author,
			})
		}
		c.JSON(http.StatusOK, AlarmListResponse{Alarms: rep})
	}
}
