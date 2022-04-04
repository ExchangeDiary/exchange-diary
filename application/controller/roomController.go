package controller

import (
	"net/http"
	"time"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/clients/google/tasks"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// RoomController handles /v1/rooms api
type RoomController interface {
	Get() gin.HandlerFunc
	GetAll() gin.HandlerFunc
	Post() gin.HandlerFunc
	Patch() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Join() gin.HandlerFunc
	Leave() gin.HandlerFunc
	GetOrders() gin.HandlerFunc
}

type roomController struct {
	roomService service.RoomService
	taskService service.TaskService
}

// NewRoomController is a roomController's constructor
func NewRoomController(rs service.RoomService, ts service.TaskService) RoomController {
	return &roomController{
		roomService: rs,
		taskService: ts,
	}
}

type responseMember struct {
	ID         uint   `json:"id"`
	NickName   string `json:"nickName"`
	ProfileURL string `json:"profileUrl"`
}

type responseRoom struct {
	ID        uint              `json:"id"`
	Name      *string           `json:"name"`
	Orders    []uint            `json:"orders"`
	Members   *[]responseMember `json:"members"`
	CreatedAt *time.Time        `json:"createdAt"`
	UpdatedAt *time.Time        `json:"updatedAt"`
}

type listResponseRoom struct {
	Rooms []responseRoom `json:"rooms"`
}

// @Summary      List rooms
// @Description  참여중인 교환일기방 리스트
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Success      200  {object}   application.JSONSuccessResponse{data=listResponseRoom,code=int,message=string}
// @Failure      400001  {object}   application.JSONBadResponse{code=int,message=string} "EntityNotFoundErr"
// @Failure      500  {object}   application.JSONServerErrResponse{code=int,message=string}
// @Router       /rooms [get]
// @Security ApiKeyAuth
func (rc *roomController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		rooms, err := rc.roomService.GetAllJoinedRooms(currentMember.ID)
		if err != nil {
			logger.Error(err.Error())
			application.FailResponse(c, application.EntityNotFoundErr, err.Error())
			return
		}
		roomsResponse := []responseRoom{}
		for _, room := range *rooms {
			members := []responseMember{}
			for _, member := range *room.Members {
				members = append(members, responseMember{
					ID:         member.ID,
					NickName:   member.Name,
					ProfileURL: member.ProfileURL,
				})
			}
			roomsResponse = append(roomsResponse, responseRoom{
				ID:        room.ID,
				Name:      &room.Name,
				Orders:    room.Orders,
				Members:   &members,
				CreatedAt: room.CreatedAt,
				UpdatedAt: room.UpdatedAt,
			})
		}
		application.Ok(c, listResponseRoom{Rooms: roomsResponse})
	}
}

type detailResponseRoom struct {
	ID            uint              `json:"id"`
	Name          *string           `json:"name"`
	Members       *[]responseMember `json:"members"`
	CreatedAt     *time.Time        `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt"`
	Theme         *string           `json:"theme,omitempty"`
	Period        uint8             `json:"period,omitempty"`
	TurnAccountID uint              `json:"turnAccountId,omitempty"`
	Code          *string           `json:"code,omitempty"`
	Hint          *string           `json:"hint,omitempty"`
	IsMaster      bool              `json:"isMaster,omitempty"`
}

// @Summary      get a room
// @Description  교환일기방 상세
// @Description  members는 교환일기방에 참여한 순서로 정렬된다.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id    path     int  true  "교환일기방 ID"  Format(uint)
// @Success      200  {object}   application.JSONSuccessResponse{data=detailResponseRoom,code=int,message=string}
// @Failure      400001  {object}   application.JSONBadResponse{code=int,message=string} "EntityNotFoundErr"
// @Failure      400002  {object}   application.JSONBadResponse{code=int,message=string} "EmptyParameterErr"
// @Failure      401003  {object}   application.JSONBadResponse{code=int,message=string} "OnlyMemberOrMasterErr"
// @Failure      500  {object}   application.JSONServerErrResponse{code=int,message=string}
// @Router       /rooms/{id} [get]
// @Security ApiKeyAuth
func (rc *roomController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			application.FailResponse(c, application.EmptyParameterErr, err.Error())
			return
		}
		room, err := rc.roomService.Get(roomID, entity.JoinedOrder)
		if err != nil {
			application.FailResponse(c, application.EntityNotFoundErr, err.Error())
			return
		}
		if !room.IsAlreadyJoined(currentMember.ID) {
			application.FailResponse(c, application.OnlyMemberOrMasterErr, "Only member or master can access")
			return
		}

		members := []responseMember{}
		for _, member := range *room.Members {
			members = append(members, responseMember{
				ID:         member.ID,
				NickName:   member.Name,
				ProfileURL: member.ProfileURL,
			})
		}
		res := detailResponseRoom{
			ID:            room.ID,
			Name:          &room.Name,
			Theme:         &room.Theme,
			Period:        room.Period,
			Members:       &members,
			TurnAccountID: room.TurnAccountID,
			CreatedAt:     room.CreatedAt,
			UpdatedAt:     room.UpdatedAt,
			IsMaster:      room.IsMaster(currentMember.ID),
		}
		application.Ok(c, res)
	}
}

type roomOrderResponse struct {
	Members       *[]responseMember `json:"members"`
	TurnAccountID uint              `json:"turnAccountId,omitempty"`
}

// @Summary      Get room orders
// @Description  교환일기방 작성 순서
// @Description  orders는 교환일기방의 다이어리 작성 순서로 정렬된다.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id    path     int  true  "교환일기방 ID"  Format(uint)
// @Success      200  {object}   roomOrderResponse
// @Success      200  {object}   application.JSONSuccessResponse{data=detailResponseRoom,code=int,message=string}
// @Failure      400001  {object}   application.JSONBadResponse{code=int,message=string} "EntityNotFoundErr"
// @Failure      400002  {object}   application.JSONBadResponse{code=int,message=string} "EmptyParameterErr"
// @Failure      401003  {object}   application.JSONBadResponse{code=int,message=string} "OnlyMemberOrMasterErr"
// @Failure      500  {object}   application.JSONServerErrResponse{code=int,message=string}
// @Router       /rooms/{id}/orders [get]
// @Security ApiKeyAuth
func (rc *roomController) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			application.FailResponse(c, application.EmptyParameterErr, err.Error())
			return
		}
		room, err := rc.roomService.Get(roomID, entity.DiaryOrder)
		if err != nil {
			application.FailResponse(c, application.EntityNotFoundErr, err.Error())
			return
		}
		if !room.IsAlreadyJoined(currentMember.ID) {
			application.FailResponse(c, application.OnlyMemberOrMasterErr, "Only member or master can access room orders")
			return
		}

		members := []responseMember{}
		for _, member := range *room.Members {
			members = append(members, responseMember{
				ID:         member.ID,
				NickName:   member.Name,
				ProfileURL: member.ProfileURL,
			})
		}

		res := roomOrderResponse{
			Members:       &members,
			TurnAccountID: room.TurnAccountID,
		}
		application.Ok(c, res)
	}
}

type postRequestRoom struct {
	Name   string `json:"name" example:"고영희방"`
	Code   string `json:"code" example:"제민욱"`
	Hint   string `json:"hint" example:"레오의 본명은?"`
	Period uint8  `json:"period" example:"5"`
	Theme  string `json:"theme" example:"1"`
}

type postResponseRoom struct {
	RoomID uint `json:"roomId"`
}

// @Summary      create a room
// @Description  교환일기방 생성
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body     postRequestRoom  true  "교환일기방 생성요청 body"
// @Success      200  {object}   postResponseRoom
// @Failure      400
// @Router       /rooms [post]
// @Security ApiKeyAuth
func (rc *roomController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		var req postRequestRoom
		if err := c.BindJSON(&req); err != nil {
			application.FailResponse(c, application.InvalidRequestBodyErr, err.Error())
			return
		}
		room, err := rc.roomService.Create(currentMember.ID, req.Name, req.Code, req.Hint, req.Theme, req.Period)
		if err != nil {
			logger.Error(err.Error())
			application.FailResponse(c, application.CannotCreateErr, err.Error())
			return
		}

		// register RoomPeriodFinCode callback task
		taskClient := tasks.GetClient()
		if _, err := rc.taskService.RegisterRoomPeriodFINTask(
			taskClient,
			application.GetCurrentURL(c),
			room.ID,
			room.TurnAccountID,
			room.DueAt,
		); err != nil {
			logger.Error(err.Error())
			application.FailResponse(c, application.TaskRoomPeriodFINCreateErr, err.Error())
			return
		}

		res := postResponseRoom{RoomID: room.ID}
		application.Ok(c, res)
	}
}

type patchRequestRoom struct {
	Code   string `json:"code,omitempty"`
	Hint   string `json:"hint,omitempty"`
	Period uint8  `json:"period,omitempty"`
	Orders []uint `json:"orders,omitempty"`
}

func (p *patchRequestRoom) ToEntity(room *entity.Room) *entity.Room {
	if p.Code != "" {
		room.Code = p.Code
	}
	if p.Hint != "" {
		room.Hint = p.Hint
	}
	if p.Period != 0 {
		// update DueAt, if period is changed
		// it will applied next turn!
		beforeDueAt := room.BeforeDueAt()
		newDueAt := beforeDueAt.Add(entity.PeriodToDuration(p.Period))
		room.Period = p.Period
		room.DueAt = &newDueAt
	}
	if p.Orders != nil {
		room.Orders = p.Orders
	}
	return room
}

func (p *patchRequestRoom) isPeriodChanged() bool {
	return p.Period != 0
}

type patchResponseRoom struct {
	RoomID uint `json:"roomId"`
}

// @Summary      update a room
// @Description  교환일기방 업데이트 (master only)
// @Description  1. 작성주기 변경 (period)
// @Description  2. 코드/힌트 변경 (code, hint)
// @Description  3. 작성순서 변경(orders) : member id를 array로 넣어주면 된다.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path 	int  true "교환일기방 ID"
// @Param        room  body 	patchRequestRoom  true "교환일기방 수정 요청 body"
// @Success      200  {object}   patchResponseRoom
// @Failure      400
// @Router       /rooms/{id} [patch]
// @Security ApiKeyAuth
func (rc *roomController) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		var req patchRequestRoom
		if err := c.BindJSON(&req); err != nil {
			application.FailResponse(c, application.InvalidRequestBodyErr, err.Error())
			return
		}

		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			application.FailResponse(c, application.EmptyParameterErr, err.Error())
			return
		}

		room, err := rc.roomService.Get(roomID, entity.Ignore)
		if err != nil {
			application.FailResponse(c, application.EntityNotFoundErr, err.Error())
			return
		}

		// TODO: 여기부터 처리할것
		if !room.IsMaster(currentMember.ID) {
			c.JSON(http.StatusUnauthorized, "Only master can patch room")
			return
		}

		_, err = rc.roomService.Update(req.ToEntity(room))
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		res := patchResponseRoom{RoomID: room.ID}
		application.Ok(c, res)
	}
}

// @Summary      (debug only) delete a room
// @Description  교환일기방 삭제
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path 	int  true "교환일기방 ID"
// @Success      204	"NO CONTENT"
// @Failure      400
// @Router       /rooms/{id} [delete]
// @Security ApiKeyAuth
func (rc *roomController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			application.FailResponse(c, application.EmptyParameterErr, err.Error())
			return
		}

		room, err := rc.roomService.Get(roomID, entity.Ignore)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if !room.IsMaster(currentMember.ID) {
			c.JSON(http.StatusUnauthorized, "Only master can delete room")
			return
		}

		err = rc.roomService.Delete(room)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		application.NoConent(c)
	}
}

type verifyRequestRoom struct {
	Code string `json:"code"`
}

// @Summary      join a room
// @Description  교환일기방 참여코드 체크 후, 교환일기방 멤버로 추가
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path 	int  true "교환일기방 ID"
// @Param        room  body 	verifyRequestRoom  true "교환일기방 참여 요청 body"
// @Success      201
// @Failure      400
// @Failure      401
// @Router       /rooms/{id}/join [post]
// @Security ApiKeyAuth
func (rc *roomController) Join() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		var req verifyRequestRoom
		if err := c.BindJSON(&req); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			application.FailResponse(c, application.EmptyParameterErr, err.Error())
			return
		}

		ok, err := rc.roomService.JoinRoom(roomID, currentMember.ID, req.Code)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if !ok {
			logger.Error(err.Error())
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}
		application.Created(c, nil)
	}
}

// @Summary      leave a room
// @Description  교환일기방 나가기
// @Description  1. 교환일기방 마스터일 경우
// @Description  2. 교환일기방 멤버일 경우
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path 	int  true "교환일기방 ID"
// @Success      204
// @Failure      400
// @Router       /rooms/{id}/leave [delete]
// @Security ApiKeyAuth
func (rc *roomController) Leave() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			application.FailResponse(c, application.EmptyParameterErr, err.Error())
			return
		}
		if err := rc.roomService.LeaveRoom(roomID, currentMember.ID); err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		application.NoConent(c)
	}
}
