package controller

import (
	"net/http"

	"github.com/exchange-diary/application"
	"github.com/exchange-diary/domain/service"
	"github.com/gin-gonic/gin"
)

// RoomController handles /v1/rooms api
type RoomController interface {
	Get() gin.HandlerFunc
	GetAll() gin.HandlerFunc
	Post() gin.HandlerFunc
	Patch() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Leave() gin.HandlerFunc
}

type roomController struct {
	roomService service.RoomService
}

// NewRoomController is a roomController's constructor
func NewRoomController(roomService service.RoomService) RoomController {
	return &roomController{roomService: roomService}
}

type postRequestRoom struct {
	Name  string `json:"name"`
	Theme string `json:"theme"`
	patchRequestRoom
}

type patchRequestRoom struct {
	Code   string `json:"code"`
	Hint   string `json:"hint"`
	Period int    `json:"period"`
	Orders int    `json:"orders"`
}

type responseRoom struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	Hint  string `json:"hint"`
	Theme string `json:"theme"`
}

// 참여중인 교환일기방 리스트
func (rc *roomController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset := application.GetLimitAndOffset(c)
		rooms, err := rc.roomService.GetAll(limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := []responseRoom{}
		for _, room := range *rooms {
			response = append(response, responseRoom{
				ID:    room.ID,
				Name:  room.Name,
				Code:  room.Code,
				Hint:  room.Hint,
				Theme: room.Theme,
			})
		}
		c.JSON(http.StatusOK, response)
	}
}

// 교환일기방 상세
func (rc *roomController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		room, err := rc.roomService.Get(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		res := responseRoom{
			ID:    room.ID,
			Name:  room.Name,
			Code:  room.Code,
			Hint:  room.Hint,
			Theme: room.Theme,
		}
		c.JSON(http.StatusCreated, res)
	}
}

// 교환일기방 생성
func (rc *roomController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req postRequestRoom
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		room, err := rc.roomService.Create(req.Name, req.Code, req.Hint, req.Theme)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		res := responseRoom{
			ID:    room.ID,
			Name:  room.Name,
			Code:  room.Code,
			Hint:  room.Hint,
			Theme: room.Theme,
		}
		c.JSON(http.StatusCreated, res)
	}
}

// 교환일기방 업데이트 (master only)
// 1. 작성주기(period)
// 2. 코드/힌트 (code, hint)
// 3. 작성순서(orders)
func (rc *roomController) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 교환일기방 삭제
func (rc *roomController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = rc.roomService.Delete(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// 교환일기방 나가기
// 1.마스터일 경우
// 2.멤버일 경우
func (rc *roomController) Leave() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
