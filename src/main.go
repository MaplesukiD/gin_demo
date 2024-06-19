package main

import (
	"gin_demo/src/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

func main() {
	//启动yml配置
	InitConfig()
	//连接数据库
	common.InitDB()

	//访问地址，处理请求
	r := gin.Default()
	r = CollectRoute(r)
	//启动服务器
	port := viper.GetString("port")
	if port != "" {
		panic(r.Run(":" + port))

	}

}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	panic(err)
}
