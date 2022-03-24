package route

import (
	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

// AlarmRoutes ...
func AlarmRoutes(router *gin.RouterGroup, controller controller.AlarmController) {
	member := router.Group("/alarms")
	{
		member.GET("/", controller.List())
	}
}
