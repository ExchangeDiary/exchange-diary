package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONSuccessResponse ...
type JSONSuccessResponse struct {
	Code    int         `json:"code" example: "200"`
	Message string      `json:"message" example: "Success"`
	Data    interface{} `json:"data"`
}

// JSONBadResponse ...
type JSONBadResponse struct {
	Code    int         `json:"code" example: "40101"`
	Message string      `json:"message" example: "Only master can access"`
	Data    interface{} `json:"data"`
}

// JSONServerErrResponse ...
type JSONServerErrResponse struct {
	Code    int         `json:"code" example: "500"`
	Message string      `json:"message" example: "Error database"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSONSuccessResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}

func FailResponse(c *gin.Context, code int, message string) {
	if code == http.StatusInternalServerError {
		c.JSON(code, JSONServerErrResponse{
			Code:    code,
			Data:    nil,
			Message: message,
		})
		return
	}
	c.JSON(code, JSONBadResponse{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}
