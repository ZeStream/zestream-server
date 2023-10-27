package controllers

import (
	"encoding/json"
	"net/http"
	"zestream-server/configs"
	"zestream-server/types"
	"zestream-server/utils"

	"github.com/gin-gonic/gin"
)

type Process struct{}

func (p Process) Video(c *gin.Context) {
	var request types.Video

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, channel, queue, ctx, _ := configs.GetRabbitMQ()

	err = utils.PublishEvent(channel, queue, *ctx, jsonBytes)

	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (p Process) Image(c *gin.Context) {
	fileName := c.Query("fileName")

	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter 'fileName' or 'fileName' provided without any extension"})
		return
	}

	// utils.Fetch()
}
