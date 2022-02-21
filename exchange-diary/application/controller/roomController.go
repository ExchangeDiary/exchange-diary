package controller

import (
	"net/http"

	"github.com/ExchangeDiary_Server/exchange-diary/domain/service"
	"github.com/gin-gonic/gin"
)


type RoomController interface {
	Get() gin.HandlerFunc
	GetAll() gin.HandlerFunc
	Post() gin.HandlerFunc
	Put() gin.HandlerFunc
	Delete() gin.HandlerFunc
}


type roomController struct {
	roomService service.RoomService
}

func NewRoomController(roomService service.RoomService) RoomController {
	return &roomController{ roomService: roomService}
}

type requestRoom struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Hint string `json:"hint"`
	Theme string `json:"theme"`
}

type responseRoom struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Hint string `json:"hint"`
	Theme string `json:"theme"`
}


func (rc *roomController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := rc.roomService.GetAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		response := []responseRoom{}
		for _, room := range *rooms {
			response = append(response, responseRoom{
				ID: room.ID,
				Name: room.Name,
			})
		}
		c.JSON(http.StatusOK, response)
	}
}

func (rc *roomController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {}
}


func (rc *roomController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requestRoom
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
			ID: room.ID,
			Name: room.Name,
			Code: room.Code,
			Hint: room.Hint,
			Theme: room.Theme,
		}
		c.JSON(http.StatusCreated, res)
	}
}

func (rc *roomController) Put() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (rc *roomController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}