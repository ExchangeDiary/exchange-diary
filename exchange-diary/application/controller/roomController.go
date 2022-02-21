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

type requestRoom struct {}

type responseRoom struct {
	ID        	int    `json:"id"`
	Name  		string `json:"name"`
	TotalMemberCount  	int `json:"totalMemberCount"`
}


func (roomController *roomController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := roomController.roomService.GetAll()
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

func (roomController *roomController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {}
}


func (roomController *roomController) Post() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (roomController *roomController) Put() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (roomController *roomController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}