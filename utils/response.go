package utils

import (
	"store/schema"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, code int, page, pageSize *int, total *int64, message string, data interface{}, isMiddleware bool) {
	response := schema.Response{}
	response.Message = message
	response.Data = data

	if page != nil {
		response.Page = *page
	}

	if pageSize != nil {
		response.PageSize = *pageSize
	}

	if total != nil {
		response.Total = *total
	}

	if isMiddleware {
		c.AbortWithStatusJSON(code, response)
		return
	} else {
		c.JSON(code, response)
		return
	}
}
