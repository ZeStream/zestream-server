package controllers

import (
	"encoding/json"
	"net/http"
	"zestream/utils"

	"github.com/gin-gonic/gin"
)

/*
	{
	  "video": {
	    "id": "string",
	    "src": "string",
	    "type": "mp4"
	  }
	}
*/
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
	message := utils.PublishMessage(utils.Config["KAFKA_URI"], string(jsonBytes), "video")
	if message != "success" {
		c.JSON(http.StatusCreated, gin.H{"status": "success"})
		return
	} else {
		c.JSON(http.StatusExpectationFailed, gin.H{"status": "fail"})
	}
}
