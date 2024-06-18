package controller

import (
	"gin_demo/src/common"
	"gin_demo/src/entity"
	"gin_demo/src/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	//密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "加密错误",
		})
		return
	}
	user := entity.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	db.Create(&user)
	//返回结果
	context.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

func Login(context *gin.Context) {
	db := common.GetDB()
	//获取参数
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

	//判断手机号是否存在
	var user entity.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
		return
	}

	//判断密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统异常",
		})
		log.Println("token生成异常:", err)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登陆成功",
		"data": gin.H{
			"token": token,
		},
	})
}

func Info(context *gin.Context) {
	user, _ := context.Get("user")

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"user": user,
		},
	})
}

// 判断数据库中不存在该手机号
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user entity.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
