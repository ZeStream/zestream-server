package controllers

import (
	"encoding/json"
	"net/http"
	"zestream-server/configs"
	"zestream-server/types"
	"zestream-server/utils"

	"github.com/gin-gonic/gin"
)

func ProcessVideo(c *gin.Context) {
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
