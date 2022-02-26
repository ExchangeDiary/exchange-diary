package route

import (
	"github.com/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

// RoomRoutes is room api handler
func RoomRoutes(incomingRoutes *gin.Engine, controller controller.RoomController) {
	rooms := incomingRoutes.Group("v1/rooms")
	{
		rooms.GET("/", controller.GetAll())
		rooms.GET("/:room_id", controller.Get())
		rooms.POST("/", controller.Post())
		rooms.PATCH("/:room_id", controller.Patch())
		rooms.DELETE("/:room_id", controller.Delete())
		rooms.POST("/:room_id/join", controller.Join())
		rooms.DELETE("/:room_id/leave", controller.Leave())
	}
}
