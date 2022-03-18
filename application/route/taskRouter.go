package route

import (
	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

// TaskRoutes to handle google cloud tasks callback endpoints
func TaskRoutes(router *gin.RouterGroup, controller controller.TaskController) {
	tasks := router.Group("/tasks")
	{
		tasks.POST("/callback", controller.HandleEvent())
		tasks.POST("/mock", controller.MockEvent())

	}
}
