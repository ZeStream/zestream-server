package controllers

import (
	"net/http"
	"zestream-server/logger"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	logger.Info(c, "ping-pong received", logger.Z{
		"input": "output",
	})
	c.String(http.StatusOK, "pong")
}
