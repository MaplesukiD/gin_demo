package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

const NAME_PREFIX = "user_"

type User struct {
	gorm.Model
	Name      string
	Telephone string
	Password  string
}

func main() {
	//创建服务
	ginServer := gin.Default()

	//连接数据库
	db := InitDB()
	//访问地址，处理请求
	ginServer.POST("/api/auth/register", func(context *gin.Context) {
		//获取参数
		name := context.PostForm("name")
		telephone := context.PostForm("telephone")
		password := context.PostForm("password")
		//数据验证
		if len(telephone) != 11 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号必须为11位",
			})
			return
		}
		if len(password) < 6 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不得小于6位",
			})
			return
		}
		if len(name) == 0 {
			name = RandomString()
		}

		log.Println(name, telephone, password)
		//判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号已存在",
			})
			return
		}
		//创建用户
		user := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&user)
		//返回结果
		context.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})

	//启动服务器
	panic(ginServer.Run())

}

// RandomString 生成随机用户名
func RandomString() string {
	return (NAME_PREFIX + uuid.New().String())[:10]
}

// 判断数据库中是否存在该手机号
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func InitDB() *gorm.DB {
	dsn := "root:1234@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败:" + err.Error())
	}
	return db
}
