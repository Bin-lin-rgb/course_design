package utils

import (
	"backend/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	response *Response
	j        *JWT
	z        = common.GetLogger()
)

// AuthUser 用户合法身份校验
func AuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.Abort()
			msg := "获取 token 失败"
			z.Error(msg)
			response.Err(ctx, http.StatusUnauthorized, msg, nil)
			return
		}
		payload, err := j.TokenParse(token)
		if err != nil {
			ctx.Abort()
			msg := "token 不合法"
			z.Error(fmt.Sprintf("%s:%s", msg, err))
			response.Err(ctx, http.StatusUnauthorized, msg, nil)
			return
		}
		z.Info(fmt.Sprintf("请求接口的用户ID：%d", payload.UserID))

		ctx.Set("userID", payload.UserID)
		token, err = j.TokenNew(payload.UserID, payload.Username) //更新token
		if err != nil {
			ctx.Abort()
			msg := "更新 token 失败"
			z.Error(fmt.Sprintf("%s:%s", msg, err))
			response.Err(ctx, http.StatusInternalServerError, msg, nil)
			return
		}
		ctx.Writer.Header().Set("Authorization", token)
		ctx.Next()

	}
}
