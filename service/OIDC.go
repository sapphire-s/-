package service

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"scutbot.cn/uniauth/database"
	"scutbot.cn/uniauth/model"
	"scutbot.cn/uniauth/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"scutbot.cn/uniauth/internal"
)

var (
	redisClient               = database.GetRedis()
	ExpiresTime time.Duration = 300000000000
)

type wellKnown struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	JwksUri                          string   `json:"jwks_uri"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                  []string `json:"scopes_supported"`
}
type Jwks struct {
	Keys []jwk.Key `json:"keys"`
}

func Authorize(context *gin.Context) {
	//获取请求中的参数
	client_id := context.Query("client_id")
	redirect_uri := context.Query("redirect_uri")
	//response_type := context.Query("response_type")
	//scope := context.Query("scope")
	state := context.Query("state")
	zap.L().Info("Authorize start")
	//暂时不管这个
	//if scope != "openid" || response_type != "code" {
	//	zap.L().Error("scope != openid || response_type != code", zap.String("scope", scope), zap.String("response_type", response_type))
	//	context.JSON(http.StatusOK, response.Result(400, "不满足oidc规范", nil))
	//	return
	//}
	client, num, err := internal.GetClientById(client_id)
	if err != nil {
		zap.L().Error("Check client_id error", zap.Error(err), zap.String("client_id", client_id))
		context.JSON(http.StatusOK, response.Result(500, "查询Client错误", nil))
		return
	}
	if num == 0 {
		zap.L().Error("Client_id error", zap.String("client_id", client_id))
		context.JSON(http.StatusOK, response.Result(400, "client_id错误", nil))
		return
	}

	//判断是否已经授权
	zap.L().Info("Client_id check")
	userinfo, err := getUserInfoFromClaims(context)
	if err != nil {
		return
	}
	audiences := strings.Split(userinfo.Audience, ",")
	var isAuthorize = false
	for _, audience := range audiences {
		if audience == client_id {
			isAuthorize = true
			break
		}
	}
	if !isAuthorize {
		zap.L().Info("client_id is not authorized", zap.String("client_id", client_id))
		context.JSON(http.StatusOK, response.Result(402, "未授权的应用，询问是否授权", client.Name))
		return
	}
	//生成一个该应用的token
	token, err := GenerateToken(userinfo, client_id)
	if err != nil {
		zap.L().Error("Generate token error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "生成token错误", nil))
		return
	}
	// 生成一个20位的随机授权码
	code := utils.RandStringRunes(20)
	//在redis中暂存授权码
	zap.L().Info("Code", zap.String("code", code))
	cmd := redisClient.SetEX(context, code, token, ExpiresTime)
	zap.L().Info("Redis cmd", zap.String("cmd", cmd.String()))
	context.JSON(http.StatusOK, response.Result(302, "Authorize重定向", redirect_uri+"?code="+code+"&state="+state))
}

func Token(context *gin.Context) {
	//获取请求中的参数
	zap.L().Info("Token start")
	grant_type := context.PostForm("grant_type")
	code := context.PostForm("code")
	if grant_type != "authorization_code" {
		zap.L().Error("Only support authorization_code", zap.String("grant_type", grant_type))
		context.JSON(http.StatusOK, response.Result(400, "仅支持授权码模式", nil))
		return
	}
	token := redisClient.Get(context, code)
	if token.Val() == "" {
		zap.L().Error("Code error", zap.String("code", code))
		context.JSON(http.StatusOK, response.Result(400, "授权码错误", nil))
		return
	}
	access_token := utils.RandStringRunes(32)
	refresh_token := utils.RandStringRunes(32)
	redisClient.Set(context, access_token, token.Val(), ExpiresTime)
	result := model.Token{
		AccessToken:  access_token,
		TokenType:    "Bearer",
		RefreshToken: refresh_token,
		ExpiresIn:    3600,
		IdToken:      token.Val(),
	}
	zap.L().Info("Return token", zap.Any("result", result))
	context.JSON(http.StatusOK, result)
}
func Info(context *gin.Context) {
	zap.L().Info("Info start")
	AccessToken := context.GetHeader("Authorization")
	zap.L().Info("AccessToken", zap.String("AccessToken", AccessToken))
	if strings.HasPrefix(AccessToken, "Bearer ") {
		// 移除 "Bearer " 部分
		AccessToken = strings.TrimPrefix(AccessToken, "Bearer ")
	} else {
		zap.L().Error("AccessToken error", zap.String("AccessToken", AccessToken))
		context.JSON(500, response.Result(500, "token格式错误", nil))
	}
	token := redisClient.Get(context, AccessToken)
	zap.L().Info("IDToken", zap.String("IDToken", token.Val()))
	j := NewJWT()
	claims, err := j.ParseToken(token.Val())
	if err != nil {
		zap.L().Error("Token parse error", zap.String("token", token.Val()), zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "token解析错误", nil))
		return
	}
	zap.L().Info("Return Info", zap.Any("claims", claims))
	context.JSON(http.StatusOK, *claims)
}

func Audience(context *gin.Context) {
	zap.L().Info("Audience start")
	aud := context.Query("aud")
	token := context.GetHeader("token")
	userinfo, err := getUserInfoFromClaims(context)
	if err != nil {
		return
	}
	audiences := strings.Split(userinfo.Audience, ",")
	for _, audience := range audiences {
		if aud == audience {
			zap.L().Info("Audience exists")
			context.JSON(http.StatusOK, response.Result(200, "已经授权的应用", token))
			return
		}
	}
	zap.L().Info("Audience not exists,refresh token")
	audiences = append(audiences, aud)
	result := strings.Join(audiences, ",")
	userinfo.Audience = result
	err = internal.UpdateUser(userinfo)
	if err != nil {
		zap.L().Error("Userinfo update error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "用户信息更新错误", token))
		return
	}
	refreshedToken, err := GenerateToken(userinfo, aud)
	if err != nil {
		zap.L().Error("refresh token error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "生成token失败", token))
		return
	}
	zap.L().Info("Audience add success,return new token", zap.String("refreshedToken", refreshedToken))
	context.JSON(http.StatusOK, response.Result(200, "授权成功", refreshedToken))
}
func WellKnown(context *gin.Context) {
	w := wellKnown{
		Issuer:                           viper.GetString("jwt.issuer"),
		AuthorizationEndpoint:            viper.GetString("wellKnown.AuthorizationEndpoint"),
		TokenEndpoint:                    viper.GetString("wellKnown.TokenEndpoint"),
		UserinfoEndpoint:                 viper.GetString("wellKnown.UserinfoEndpoint"),
		JwksUri:                          viper.GetString("wellKnown.JwksUri"),
		SubjectTypesSupported:            viper.GetStringSlice("wellKnown.SubjectTypesSupported"),
		IdTokenSigningAlgValuesSupported: viper.GetStringSlice("wellKnown.IdTokenSigningAlgValuesSupported"),
		ResponseTypesSupported:           viper.GetStringSlice("wellKnown.ResponseTypesSupported"),
		ScopesSupported:                  viper.GetStringSlice("wellKnown.ScopesSupported"),
	}
	zap.L().Info("WellKnown Send", zap.Any("wellKnown", w))
	context.JSON(http.StatusOK, w)
}
func Jwk(context *gin.Context) {
	keys := make([]jwk.Key, 0)
	key, _ := GenerateJWK()
	keys = append(keys, key)
	jsonWebKey := Jwks{
		keys,
	}
	zap.L().Info("Jwk Send", zap.Any("jwk", jsonWebKey))
	context.JSON(http.StatusOK, jsonWebKey)
}
