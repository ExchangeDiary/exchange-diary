package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONSuccessResponse ...
type JSONSuccessResponse struct {
	Code    uint        `json:"code" example:"200"`
	Message string      `json:"message" example:"Ok"`
	Data    interface{} `json:"data"`
}

// JSONBadResponse ...
type JSONBadResponse struct {
	Code    int    `json:"code" example:"40001"`
	Message string `json:"message" example:"4xx errors"`
}

// JSONServerErrResponse ...
type JSONServerErrResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Error database"`
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSONSuccessResponse{
		Code:    http.StatusOK,
		Message: "Ok",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSONSuccessResponse{
		Code:    http.StatusCreated,
		Message: "Created",
		Data:    data,
	})
}

func NoConent(c *gin.Context) {
	c.JSON(http.StatusOK, JSONSuccessResponse{
		Code:    http.StatusNoContent,
		Message: "No Conent",
		Data:    c.Request.TLS.NegotiatedProtocolIsMutual,
	})
}

func FailResponse(c *gin.Context, code int, message string) {
	if code == http.StatusInternalServerError {
		c.JSON(code, JSONServerErrResponse{
			Code:    code,
			Message: message,
		})
		return
	}
	c.JSON(code, JSONBadResponse{
		Code:    code,
		Message: message,
	})
}
