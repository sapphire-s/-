package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"scutbot.cn/uniauth/middleware"
	"scutbot.cn/uniauth/service"
)

func InitRouter() {
	r := gin.Default()
	// 注册
	r.POST("/register", service.Register)
	// 登录
	r.POST("/login", service.Login)
	// 修改密码
	r.POST("/ChangePasswd", middleware.CheckToken(), service.ChangePasswd)
	//修改头像
	r.POST("/Avatar", middleware.CheckToken(), service.ChangeAvatar)
	//获取头像
	r.GET("/Avatar", service.GetAvatar)
	// OIDC
	oidc := r.Group("/api")
	{
		oidc.POST("/authorize", middleware.CheckToken(), service.Authorize)
		oidc.GET("/authorize", middleware.CheckToken(), service.Authorize)
		oidc.POST("/token", service.Token)
		oidc.GET("/info", service.Info)
		oidc.GET("/.well-known/openid-configuration", service.WellKnown)
		oidc.POST("/audience", middleware.CheckToken(), service.Audience)
		oidc.GET("/jwks", service.Jwk)
	}
	// 飞书相关
	feishu := r.Group("/feishu")
	{
		// 飞书登录
		feishu.GET("/login", service.FeishuLogin)
		// 飞书登录回调接口
		feishu.GET("/callback", service.Callback)
		// 绑定飞书ID和uuid
		feishu.POST("/bind", middleware.CheckToken(), service.Bind)
	}
	//身份
	identity := r.Group("/id")
	{
		identity.POST("/", middleware.CheckToken(), service.Apply)
		identity.GET("/", middleware.CheckToken(), service.GetIdentity)
		identity.DELETE("/", middleware.CheckToken(), service.Withdraw)
		identity.PUT("/", middleware.CheckToken(), middleware.CheckPermission("All"), service.Approve)
		identity.GET("/pend", middleware.CheckToken(), service.Pending)
	}
	r.POST("/casbin/add", service.AddPermission)
	// Auto
	r.GET("/auto", service.AutoMigrate)
	// 测试
	r.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "测试成功",
		})
	})
	r.POST("/casbin/test", middleware.CheckToken(), middleware.CheckPermission("24软开组"), func(context *gin.Context) {
		context.JSON(200, "有权限的")
	})
	r.GET("/casbin/test", middleware.CheckToken(), middleware.CheckPermission("24软开组"), func(context *gin.Context) {
		context.JSON(200, "有权限的")
	})
	r.Run(":" + viper.GetString("server.addr"))
}
