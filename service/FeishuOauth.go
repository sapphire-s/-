package service

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"scutbot.cn/uniauth/internal"
	"scutbot.cn/uniauth/lark"
	"scutbot.cn/uniauth/model"
)

var (
	BaseURL     = viper.Get("feishu.base_url").(string)
	AppId       = "app_id=" + viper.Get("feishu.app_id").(string)
	RedirectURI = "redirect_uri=" + viper.Get("feishu.redirect_uri").(string)
	BindURI     = viper.Get("feishu.bind_uri").(string)
	loginURI    = viper.Get("feishu.login_uri").(string)
)

func FeishuLogin(context *gin.Context) {
	context.Redirect(http.StatusFound, BaseURL+AppId+"&"+RedirectURI)
}

func Callback(context *gin.Context) {
	// 获取预授权码
	zap.L().Info("Code start")
	code := context.Query("code")
	zap.L().Info("code", zap.String("code", code))
	// 获取user_access_token
	userAccessToken, err := lark.GetUserToken(code)
	if err != nil {
		zap.L().Error("获取飞书user_access_token失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, response.Result(500, "获取飞书user_access_token失败", err))
		return
	}
	zap.L().Info("user_access_token", zap.String("user_access_token", userAccessToken))
	// 获取UserInfo
	userInfo, err := lark.GetUserInfo(userAccessToken)
	if err != nil {
		zap.L().Error("获取飞书用户信息失败", zap.Error(err))
		context.JSON(http.StatusInternalServerError, response.Result(500, "获取飞书用户信息失败", err))
		return
	}
	feishuUserInfo := internal.NewFeishuUserInfo(userInfo)
	_, num, _ := internal.GetFeishuMapByUnionId(feishuUserInfo.UnionId)
	// 新用户
	if num == 0 {
		zap.L().Info("新用户,开始绑定")
		// 绑定统一认证账户
		context.Redirect(http.StatusFound, BindURI+feishuUserInfo.UnionId)
	} else {
		// 已经绑定的用户,返回飞书unionId对应的token
		uuid, err := internal.GetUuid(feishuUserInfo.UnionId)
		if err != nil {
			zap.L().Error("获取uuid失败", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "获取uuid失败", err))
			return
		}
		user, num, _ := internal.GetUserByID(uuid)
		if num == 0 {
			zap.L().Error("用户不存在", zap.String("uuid", uuid))
			context.JSON(http.StatusOK, response.Result(500, "用户不存在", nil))
			return
		}
		//生成token
		token, err := GenerateToken(user, user.Audience)
		if err != nil {
			zap.L().Error("Token generate fail", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "生成token失败", err))
			return
		}
		zap.L().Info("飞书登录成功")
		context.Redirect(http.StatusFound, loginURI+token)
		return
	}
}

func Bind(context *gin.Context) {
	//获取请求中的信息
	unionId := context.Query("unionId")
	claim, _ := context.Get("claims")
	uuid := claim.(*model.JwtClaims).Uuid
	Map := model.FeishuMap{
		Uuid:    uuid,
		UnionId: unionId,
	}
	//添加映射
	err := internal.AddFeishuMap(&Map)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	zap.L().Info("绑定成功")
	context.JSON(http.StatusOK, response.Result(200, "绑定成功", nil))
}
