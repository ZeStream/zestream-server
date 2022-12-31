package helpers

import "github.com/gin-gonic/gin"

func errorSchema(err error, msg string, details string, statusCode uint8) gin.H {
	return gin.H{
		"error":      err.Error(),
		"message":    msg,
		"detail":     details,
		"statusCode": statusCode, // not sure if this should be done
	}
}
