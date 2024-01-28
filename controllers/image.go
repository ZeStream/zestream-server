package controllers

import (
	"zestream-server/utils"

	"github.com/gin-gonic/gin"
)

type Image struct{}

func (*Image) Get(c *gin.Context) {
	subpath := c.Param("subpath")
	d := utils.GetQueryHash(c.Request.URL.RawQuery)
	c.String(200, subpath+":"+d)
}
