package route

import (
	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

// FileRoutes ...
func FileRoutes(router *gin.RouterGroup, controller controller.FileController) {
	files := router.Group("/rooms/:room_id/files")
	{
		// files.GET("/:file_uuid", controller.Get())
		files.POST("/", controller.Post())
	}
}
