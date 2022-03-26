package controller

import (
	"fmt"
	"net/http"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// TaskController ...
type TaskController interface {
	HandleEvent() gin.HandlerFunc
	MockEvent() gin.HandlerFunc
}

type taskController struct {
	taskService   service.TaskService
	memberService service.MemberService
}

// NewTaskController ...
func NewTaskController(ts service.TaskService, ms service.MemberService) TaskController {
	return &taskController{
		taskService:   ts,
		memberService: ms,
	}
}

type taskRequest struct {
	RoomID uint        `json:"room_id"`
	Email  string      `json:"email"`
	Code   vo.TaskCode `json:"code" enums:"ROOM_PERIOD_FIN,MEMBER_ON_DUTY,MEMBER_BEFORE_1HR,MEMBER_BEFORE_4HR,MEMBER_POSTED_DIARY"`
}

// @Summary      Handle Event Task
// @Description	 google cloud task에 예약 해두었던, task들을 스케쥴된 일정시간이 지난뒤, 처리해주는 callback handler api endpoint.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body     taskRequest  true  "Task Http Body"
// @Success      202  "Accepted. It ignore request, because error occured. 정상 처리(204)와 차이점을 두기 위해서 202로 처리함"
// @Success      204 "successfully finished callback task."
// @Router       /tasks/callback [post]
func (tc *taskController) HandleEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req taskRequest
		if err := c.BindJSON(&req); err != nil {
			logger.Error("POST /task/callback endpoint gets invalid taskRequest json body.")
			// google cloud task에서 쏘기 때문에 204 처리해주어야 한다.
			c.JSON(http.StatusAccepted, err.Error())
			return
		}
		_, err := tc.memberService.GetByEmail(req.Email)
		if err != nil {
			logger.Error("Callback task errors. " + req.Email + " not found. (maybe user is already unsigned)")
			c.JSON(http.StatusAccepted, err.Error())
			return
		}

		err = tc.doTask(req, application.GetCurrentURL(c))
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusAccepted, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// @Summary      Mock Handle Event Task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body     taskRequest  true  "Task Http Body"
// @Success      202  "Accepted. It ignore request, because error occured. 정상 처리(204)와 차이점을 두기 위해서 202로 처리함"
// @Success      204 "successfully finished callback task."
// @Router       /tasks/mock [post]
func (tc *taskController) MockEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req taskRequest
		if err := c.BindJSON(&req); err != nil {
			logger.Error("POST /task/callback endpoint gets invalid taskRequest json body.")
			// google cloud task에서 쏘기 때문에 204 처리해주어야 한다.
			c.JSON(http.StatusAccepted, err.Error())
			return
		}
		_, err := tc.memberService.GetByEmail(req.Email)
		if err != nil {
			logger.Error("Callback task errors. " + req.Email + " not found. (maybe user is already unsigned)")
			c.JSON(http.StatusAccepted, err.Error())
			return
		}

		err = tc.doMockTask(req, application.GetCurrentURL(c))
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusAccepted, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func (tc *taskController) doMockTask(dto taskRequest, baseURL string) (err error) {
	if err = tc.taskService.DoMemberOnDutyTask(dto.RoomID, dto.Email); err != nil {
		return
	}
	if err = tc.taskService.DoMemberBeforeTask(dto.RoomID, dto.Email, vo.MemberBefore1HRCode); err != nil {
		return
	}
	if err = tc.taskService.DoMemberBeforeTask(dto.RoomID, dto.Email, vo.MemberBefore4HRCode); err != nil {
		return
	}
	if err = tc.taskService.DoMemberPostedDiaryTask(dto.RoomID, baseURL); err != nil {
		return
	}
	return
}

func (tc *taskController) doTask(dto taskRequest, baseURL string) (err error) {
	switch dto.Code {
	case vo.RoomPeriodFinCode:
		err = tc.taskService.DoRoomPeriodFINTask(dto.RoomID, baseURL)
	case vo.MemberOnDutyCode:
		err = tc.taskService.DoMemberOnDutyTask(dto.RoomID, dto.Email)
	case vo.MemberBefore1HRCode:
		err = tc.taskService.DoMemberBeforeTask(dto.RoomID, dto.Email, vo.MemberBefore1HRCode)
	case vo.MemberBefore4HRCode:
		err = tc.taskService.DoMemberBeforeTask(dto.RoomID, dto.Email, vo.MemberBefore4HRCode)
	case vo.MemberPostedDiaryCode:
		err = tc.taskService.DoMemberPostedDiaryTask(dto.RoomID, baseURL)
	default:
		err = fmt.Errorf("Not registered task code. [ " + string(dto.Code) + " ]")
	}
	return
}
