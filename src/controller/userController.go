package controller

import (
	"gin_demo/src/common"
	"gin_demo/src/entity"
	"gin_demo/src/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	db := common.GetDB()
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
		name = utils.RandomString()
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
	user := entity.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	db.Create(&user)
	//返回结果
	context.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

// 判断数据库中是否存在该手机号
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user entity.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
