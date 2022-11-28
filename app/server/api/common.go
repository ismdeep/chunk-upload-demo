package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fail(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  err.Error(),
	})
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}
