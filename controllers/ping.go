package controllers

import (
	"errors"
	"net/http"
	"zestream/constants"
	"zestream/helpers"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(
		http.StatusBadRequest,
		helpers.ResponseSchema(
			helpers.CreateErrorResponse(constants.StatusTexts[http.StatusBadRequest], "Missing some stuff", errors.New("This is an error msg")),
			gin.H{},
			http.StatusBadRequest))
}
