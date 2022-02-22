package application

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 10
	defaultOffset = 0
)

func GetLimitAndOffset(c *gin.Context) (int, int){
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = defaultOffset
	}
	return limit, offset
}