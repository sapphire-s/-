package middleware

import (
	"go.uber.org/zap"
	"net/url"

	"github.com/gin-gonic/gin"

	"scutbot.cn/uniauth/model"
	"scutbot.cn/uniauth/service"
)

var response = &model.Response{}

func CheckToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("token")
		zap.L().Info("CheckToken start", zap.String("token", token))
		if token == "" {
			zap.L().Info("token is empty", zap.String("token", token))
			reLogin(context)
			context.Abort()
			return
		}
		j := service.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			zap.L().Info("token is invalid", zap.Error(err))
			context.JSON(200, response.Result(401, "token is invalid", err))
			context.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		context.Set("claims", claims)
	}
}

func reLogin(context *gin.Context) {
	//重新登录
	originUrl := context.Request.URL
	newUrl, _ := url.Parse("https://auth.scutbot.icu/authorize")
	newUrl.RawQuery = originUrl.RawQuery
	zap.L().Info("redirectToLogin", zap.String("url", newUrl.String()))
	context.Redirect(302, newUrl.String())
}
