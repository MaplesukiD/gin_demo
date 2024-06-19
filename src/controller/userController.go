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
		common.Response(context, http.StatusUnprocessableEntity, 422, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		common.Response(context, http.StatusUnprocessableEntity, 422, "密码不得小于6位")
		return
	}
	if len(name) == 0 {
		name = utils.RandomString()
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		common.Response(context, http.StatusUnprocessableEntity, 422, "手机号已存在")
		return
	}
	//创建用户
	//密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		common.Response(context, http.StatusInternalServerError, 500, "加密错误")
		return
	}
	user := entity.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	db.Create(&user)
	//返回结果
	common.Response(context, http.StatusOK, 200, "注册成功")
}

func Login(context *gin.Context) {
	db := common.GetDB()
	//获取参数
	telephone := context.PostForm("telephone")
	password := context.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		common.Response(context, http.StatusUnprocessableEntity, 422, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		common.Response(context, http.StatusUnprocessableEntity, 422, "密码不得小于6位")
		return
	}

	//判断手机号是否存在
	var user entity.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		common.Response(context, http.StatusUnprocessableEntity, 422, "用户不存在")
		return
	}

	//判断密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		common.Response(context, http.StatusBadRequest, 400, "密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		common.Response(context, http.StatusInternalServerError, 500, "系统异常")
		log.Println("token生成异常:", err)
		return
	}
	common.ResponseWithData(context, http.StatusOK, 200, gin.H{"token": token}, "登陆成功")
}

func Info(context *gin.Context) {
	user, _ := context.Get("user")
	common.ResponseWithData(context, http.StatusOK, 200, gin.H{"user": user}, "查询成功")
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
