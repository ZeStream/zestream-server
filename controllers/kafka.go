package controllers

import (
	"encoding/json"
	"net/http"
	"zestream-server/configs"
	"zestream-server/utils"

	"github.com/gin-gonic/gin"
)

type Video struct {
	ID   string `json:"id"`
	Src  string `json:"src"`
	Type string `json:"type"`
}

type Body struct {
	Video Video `json:"video"`
}

func PublishMessage(c *gin.Context) {
	var request Body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := utils.PublishMessage(c, configs.EnvVar[configs.KAFKA_URI], string(jsonBytes), "video")
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": message})
}
