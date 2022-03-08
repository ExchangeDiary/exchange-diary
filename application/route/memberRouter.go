package route

import (
	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

func MemberRoutes(router *gin.RouterGroup, controller controller.MemberController) {
	member := router.Group("/member")
	{
		member.GET("/:email", controller.Get())
		member.POST("/", controller.Post())
		member.PATCH("/", controller.Patch())
		member.DELETE("/:email", controller.Delete())
	}
}
