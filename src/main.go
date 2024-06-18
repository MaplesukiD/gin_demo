package main

import (
	"gin_demo/src/common"
	"github.com/gin-gonic/gin"
)

func main() {
	//连接数据库
	common.InitDB()

	//访问地址，处理请求
	r := gin.Default()
	r = CollectRoute(r)
	//启动服务器
	panic(r.Run())

}
