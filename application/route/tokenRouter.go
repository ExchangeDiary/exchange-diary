package route

import (
	"github.com/ExchangeDiary/exchange-diary/application/controller"
	"github.com/gin-gonic/gin"
)

// TokenRoutes ...
func TokenRoutes(router *gin.RouterGroup, controller controller.TokenController) {
	token := router.Group("/token")
	{
		token.GET("", controller.GetToken())
		token.GET("/refresh", controller.RefreshAccessToken())
	}
}
