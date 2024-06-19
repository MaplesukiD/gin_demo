package common

import (
	"github.com/gin-gonic/gin"
)

// ResponseWithData 带Data的Response返回
func ResponseWithData(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

// Response 通用Response返回
func Response(ctx *gin.Context, httpStatus int, code int, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
	})
}
