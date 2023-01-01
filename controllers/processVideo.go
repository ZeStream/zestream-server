package controllers

import (
	"fmt"
	"zestream/constants"
	"zestream/helpers"

	"github.com/gin-gonic/gin"
)

func ProcessVideo(c *gin.Context) {
	fmt.Println("Do Something")
	c.JSON(constants.StatusOK, helpers.ResponseSchema(gin.H{}, helpers.CreateSuccessResponse("signedURL", "2415"), constants.StatusOK))
}
