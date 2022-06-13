package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(c *gin.Context) {
	name := c.Query("name")
	msg := "Hello " + name
	c.JSON(http.StatusOK, map[string]interface{}{
		"Code": 200,
		"Msg":  msg,
	})
}
