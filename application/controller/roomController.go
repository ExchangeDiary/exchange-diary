package controller

import (
	"net/http"
	"time"

	"github.com/ExchangeDiary/exchange-diary/application"
	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/service"
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
}

type roomController struct {
	roomService service.RoomService
}

// NewRoomController is a roomController's constructor
func NewRoomController(roomService service.RoomService) RoomController {
	return &roomController{roomService: roomService}
}

// TODO: move to account
type responseMember struct {
	ID         uint   `json:"id"`
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
// @Param        limit    query     uint  false  "page size"  Format(uint)
// @Param        offset    query    uint  false  "page offset"  Format(uint)
// @Success      200  {object}   listResponseRoom
// @Failure      400
// @Router       /rooms [get]
func (rc *roomController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		limit, offset := application.GetLimitAndOffset(c)
		rooms, err := rc.roomService.GetAllJoinedRooms(currentMember.ID, limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomsResponse := []responseRoom{}
		for _, room := range *rooms {
			members := []responseMember{}
			for _, member := range *room.Members {
				members = append(members, responseMember{
					ID:         member.ID,
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
		c.JSON(http.StatusOK, listResponseRoom{Rooms: roomsResponse})
	}
}

type detailResponseRoom struct {
	ID              uint              `json:"id"`
	Name            *string           `json:"name"`
	Members         *[]responseMember `json:"members"`
	CreatedAt       *time.Time        `json:"createdAt"`
	UpdatedAt       *time.Time        `json:"updatedAt"`
	Theme           *string           `json:"theme,omitempty"`
	Period          uint8             `json:"period,omitempty"`
	TurnAccountID   uint              `json:"turnAccountId,omitempty"`
	TurnAccountName *string           `json:"turnAccountName,omitempty"`
	Code            *string           `json:"code,omitempty"`
	Hint            *string           `json:"hint,omitempty"`
	IsMaster        bool              `json:"isMaster,omitempty"`
}

// @Summary      get a room
// @Description  교환일기방 상세
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id    path     int  true  "교환일기방 ID"  Format(uint)
// @Success      200  {object}   detailResponseRoom
// @Failure      400
// @Router       /rooms/{id} [get]
func (rc *roomController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
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

		members := []responseMember{}
		for _, member := range *room.Members {
			members = append(members, responseMember{
				ID:         member.ID,
				ProfileURL: member.ProfileURL,
			})
		}
		turnAccountName := "MOCK 어카운트 이름"
		res := detailResponseRoom{
			ID:              room.ID,
			Name:            &room.Name,
			Theme:           &room.Theme,
			Period:          room.Period,
			Members:         &members,
			TurnAccountID:   room.TurnAccountID,
			TurnAccountName: &turnAccountName,
			CreatedAt:       room.CreatedAt,
			UpdatedAt:       room.UpdatedAt,
			IsMaster:        room.IsMaster(currentMember.ID),
		}
		c.JSON(http.StatusOK, res)
	}
}

type postRequestRoom struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Hint   string `json:"hint"`
	Period uint8  `json:"period"`
	Theme  string `json:"theme"`
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
func (rc *roomController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		var req postRequestRoom
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		room, err := rc.roomService.Create(currentMember.ID, req.Name, req.Code, req.Hint, req.Theme, req.Period)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		res := postResponseRoom{RoomID: room.ID}
		c.JSON(http.StatusOK, res)
	}
}

type patchRequestRoom struct {
	Code    string `json:"code,omitempty"`
	Hint    string `json:"hint,omitempty"`
	Period  uint8  `json:"period,omitempty"`
	Members []uint `json:"members,omitempty"`
}

func (p *patchRequestRoom) ToEntity(room *entity.Room) *entity.Room {
	if p.Code != "" {
		room.Code = p.Code
	}
	if p.Hint != "" {
		room.Hint = p.Hint
	}
	if p.Period != 0 {
		room.Period = p.Period
	}
	if p.Members != nil {
		room.Orders = p.Members
	}
	return room
}

type patchResponseRoom struct {
	RoomID uint `json:"roomId"`
}

// @Summary      update a room
// @Description  교환일기방 업데이트 (master only)
// @Description  1. 작성주기 변경 (period)
// @Description  2. 코드/힌트 변경 (code, hint)
// @Description  3. 작성순서 변경(orders)
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path 	int  true "교환일기방 ID"
// @Param        room  body 	patchRequestRoom  true "교환일기방 수정 요청 body"
// @Success      200  {object}   patchResponseRoom
// @Failure      400
// @Router       /rooms/{id} [patch]
func (rc *roomController) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		var req patchRequestRoom
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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

		if !room.IsMaster(currentMember.ID) {
			c.JSON(http.StatusUnauthorized, "Only master can patch room")
			return
		}

		_, err = rc.roomService.Update(req.ToEntity(room))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		res := patchResponseRoom{RoomID: room.ID}
		c.JSON(http.StatusOK, res)
	}
}

// @Summary      (debug only) delete a room
// @Description  교환일기방 삭제
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id  path 	int  true "교환일기방 ID"
// @Success      204
// @Failure      400
// @Router       /rooms/{id} [delete]
func (rc *roomController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
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

		if !room.IsMaster(currentMember.ID) {
			c.JSON(http.StatusUnauthorized, "Only master can delete room")
			return
		}

		err = rc.roomService.Delete(room)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
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
func (rc *roomController) Join() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		var req verifyRequestRoom
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ok, err := rc.roomService.JoinRoom(roomID, currentMember.ID, req.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}
		c.Status(http.StatusCreated)
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
func (rc *roomController) Leave() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMember := c.MustGet(application.CurrentMemberKey).(application.CurrentMemberDTO)
		roomID, err := application.ParseUint(c.Param("room_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := rc.roomService.LeaveRoom(roomID, currentMember.ID); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}
