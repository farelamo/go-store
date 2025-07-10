package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessWithPaginate(c *gin.Context, msg string, totalData int64, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": msg,
		"total":   totalData,
		"data":    data,
	})
}

func Success(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": msg,
		"data":    data,
	})
}
