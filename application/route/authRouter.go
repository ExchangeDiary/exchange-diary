package route

import (
	"github.com/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, controller controller.AuthController) {
	redirectLogin := router.Group("/login")
	{
		redirectLogin.GET("/:auth_type", controller.Redirect())
	}
	auth := router.Group("/authentication")
	{
		auth.GET("/login/:auth_type", controller.Login())
		auth.GET("/authenticated", controller.Authenticate())
	}
}
