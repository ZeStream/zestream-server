package helpers

import "github.com/gin-gonic/gin"

func CreateSuccessResponse(preSignedURL, vidID string) gin.H {
	return gin.H{
		"videoID":      vidID,
		"preSignedURL": preSignedURL,
	}
}

func CreateErrorResponse(msg, details string, err error) gin.H {
	return gin.H{
		"message": msg,
		"details": details,
		"error":   err.Error(),
	}
}

func ResponseSchema(err gin.H, successRes gin.H, statusCode uint8) gin.H {
	return gin.H{
		"statusCode": statusCode,
		"data":       successRes,
		"error":      err,
	}
}

// use it like this
// ErrorSchema(CreateErrorResponse(---), {}, constant.StatusCode)
