package main

import (
	"gin_demo/src/controller"
	"gin_demo/src/middlewares"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middlewares.AuthMiddleware(), controller.Info)
	return r
}
