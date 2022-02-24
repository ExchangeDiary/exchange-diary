package application

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

// GetLimitAndOffset parse limit and offset from http context and convert them to uint type.
func GetLimitAndOffset(c *gin.Context) (uint, uint) {
	limit, err := ParseUint(c.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}
	offset, err := ParseUint(c.Query("offset"))
	if err != nil {
		offset = defaultOffset
	}
	return limit, offset
}

// ParseUint parse string to uint
func ParseUint(str string) (uint, error) {
	val, err := strconv.ParseUint(str, 10, 64) // uint64
	return uint(val), err
}
