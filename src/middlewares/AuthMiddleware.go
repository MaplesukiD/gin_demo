package middlewares

import (
	"gin_demo/src/common"
	"gin_demo/src/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		//获取header
		tokenString := context.GetHeader("Authorization")

		//验证格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}

		//获取claims中的信息,查找对应id信息
		userId := claims.UserId
		db := common.GetDB()
		var user entity.User
		db.First(&user, userId)

		//用户不存在
		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}

		context.Set("user", user)
		context.Next()
	}
}
