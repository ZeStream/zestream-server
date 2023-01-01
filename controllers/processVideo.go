package controllers

import (
	"fmt"
	"net/http"
	"zestream/helpers"

	"github.com/gin-gonic/gin"
)

func ProcessVideo(c *gin.Context) {
	fmt.Println("Do Something")
	c.JSON(
		http.StatusOK,
		helpers.ResponseSchema(gin.H{}, helpers.CreateSuccessResponse("signedURL", "2415"),
			http.StatusOK))
}
