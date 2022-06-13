package utils

import (
	"dousheng/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// response封装
//普通成功返回
func Success(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Code": common.OK,
		"Data": v,
	})
}
