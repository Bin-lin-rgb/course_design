package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SuccessCode = 0
	ErrCode     = 500
)

type Response struct {
}

func (r *Response) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    SuccessCode,
		"message": "success",
		"result":  data,
	})
	return
}

func (r *Response) Err(c *gin.Context, httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, map[string]interface{}{
		"code":    ErrCode,
		"message": msg,
		"result":  data,
	})
	return
}
