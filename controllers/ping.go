package controllers

import (
	"zestream/constants"
	"zestream/helpers"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(
		constants.StatusOK,
		helpers.ResponseSchema(
			gin.H{},
			helpers.CreateSuccessResponse("presignedURL", "2415"),
			constants.StatusOK))
}
