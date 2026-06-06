package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getUserContext(c *gin.Context) (uint, string, error) {
	userIDRaw, ok := c.Get("userID")
	if !ok {
		return 0, "", errors.New("unauthorized")
	}
	userRoleRaw, ok := c.Get("userRole")
	if !ok {
		return 0, "", errors.New("unauthorized")
	}
	userID, err := strconv.ParseUint(userIDRaw.(string), 10, 64)
	if err != nil {
		return 0, "", err
	}
	return uint(userID), strings.ToLower(userRoleRaw.(string)), nil
}
