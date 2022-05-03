package route

import (
	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

// AuthRoutes ...
func AuthRoutes(router *gin.RouterGroup, controller controller.AuthController) {
	login := router.Group("/login")
	{
		login.GET("/:auth_type", controller.Login())
	}
	auth := router.Group("/authentication")
	{
		auth.GET("/authenticated", controller.Authenticate())
		auth.POST("/mock", controller.MockRegister())
	}
}
