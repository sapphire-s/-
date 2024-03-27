package service

import (
	"github.com/gin-gonic/gin"
	"scutbot.cn/uniauth/bootstrap"
	"scutbot.cn/uniauth/model"
)

var db = bootstrap.GetDatabase()

func AutoMigrate(context *gin.Context) {
	db.AutoMigrate(&model.UserInfo{})
	db.AutoMigrate(&model.UserIdentity{})
	db.AutoMigrate(&model.Client{})
	db.AutoMigrate(&model.FeishuUserInfo{})
	db.AutoMigrate(&model.FeishuMap{})
	context.JSON(200, gin.H{
		"message": "成功了...吗？",
	})
}
