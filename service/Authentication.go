package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"scutbot.cn/uniauth/internal"
	"scutbot.cn/uniauth/model"
	"scutbot.cn/uniauth/utils"
)

var response = &model.Response{}

func Register(context *gin.Context) {
	// 接收数据
	userinfo := model.UserInfo{}
	err := context.ShouldBindJSON(&userinfo)
	if err != nil {
		zap.L().Error("Parse info error", zap.Error(err), zap.Any("userinfo", userinfo))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userinfo.Passwd), bcrypt.DefaultCost)
	//查看是否有重复邮箱
	_, num, _ := internal.GetUserByEmail(userinfo.Email)
	if num != 0 {
		zap.L().Error("Email already exists")
		context.JSON(http.StatusOK, response.Result(400, "已经存在的邮箱", nil))
		return
	}
	// 写入数据库
	user := internal.NewUserInfo(userinfo.Name, userinfo.Email, string(hashedPassword))
	//准备默认头像
	sourcePath := "./avatar/Avatar_default"
	destinationPath := "./avatar/Avatar_" + user.Uuid
	err = utils.CopyFile(sourcePath, destinationPath)
	if err != nil {
		zap.L().Error("Copy file error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "Copy file error", err))
		return
	}
	err = internal.AddUser(user)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	} else {
		token, err := GenerateToken(user, "")
		if err != nil {
			zap.L().Error("Token generate fail", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "生成token失败", err))
			return
		}
		zap.L().Info("Register success")
		context.JSON(http.StatusOK, response.Result(200, "注册成功", token))
		return
	}
}

func Login(context *gin.Context) {
	// 接收数据
	zap.L().Info("Login start")
	userinfo := model.UserInfo{}
	err := context.ShouldBindJSON(&userinfo)
	if err != nil {
		zap.L().Error("Parse info error", zap.Error(err), zap.Any("userinfo", userinfo))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	// 获取真实密码
	user, num, err := internal.GetUserByEmail(userinfo.Email)
	if num == 0 {
		zap.L().Error("Email does not exist", zap.Error(err), zap.Any("userinfo", userinfo))
		context.JSON(http.StatusOK, response.Result(400, "邮箱不存在", nil))
		return
	}
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	// 密码比对
	err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(userinfo.Passwd))
	if err != nil {
		zap.L().Error("Password or Email error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(400, "邮箱或密码错误", nil))
		return
	} else {
		token, err := GenerateToken(user, userinfo.Audience)
		if err != nil {
			zap.L().Error("Token generate fail", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "生成token失败", err))
			return
		}
		zap.L().Info("Login success")
		context.JSON(http.StatusOK, response.Result(200, "登录成功", token))
		return
	}
}

func ChangePasswd(context *gin.Context) {
	// 接收数据
	zap.L().Info("ChangePasswd start")
	userinfo := model.UserInfo{}
	err := context.ShouldBindJSON(&userinfo)
	if err != nil {
		zap.L().Error("Parse info error", zap.Error(err), zap.Any("userinfo", userinfo))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	// 新密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userinfo.Passwd), bcrypt.DefaultCost)
	// 获取用户信息
	user, err := getUserInfoFromClaims(context)
	if err != nil {
		return
	}
	user, num, err := internal.GetUserByID(user.Uuid)
	if num == 0 {
		zap.L().Error("Uuid error", zap.Any("userinfo", userinfo))
		context.JSON(http.StatusOK, response.Result(400, "用户不存在", nil))
		return
	}
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	// 更新用户信息
	user.Passwd = string(hashedPassword)
	err = internal.UpdateUser(user)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	} else {
		zap.L().Info("ChangePasswd success")
		context.JSON(http.StatusOK, response.Result(200, "密码更改成功", nil))
		return
	}
}

func ChangeAvatar(context *gin.Context) {
	zap.L().Info("ChangeAvatar start")
	avatar, _ := context.FormFile("avatar")
	userinfo, err := getUserInfoFromClaims(context)
	if err != nil {
		return
	}
	dst := "./avatar/" + "Avatar_" + userinfo.Uuid
	// 上传文件至指定的完整文件路径
	err = context.SaveUploadedFile(avatar, dst)
	if err != nil {
		zap.L().Error("SaveUploadedFile error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "头像上传失败", err))
		return
	}
	context.JSON(http.StatusOK, response.Result(200, "头像上传成功", nil))
}
func GetAvatar(context *gin.Context) {
	zap.L().Info("GetAvatar start")
	uuid := context.Query("avatar")
	dst := "./avatar/" + "Avatar_" + uuid
	context.File(dst)
}

func getUserInfoFromClaims(context *gin.Context) (*model.UserInfo, error) {
	claims, _ := context.Get("claims")
	userinfo, num, err := internal.GetUserByID(claims.(*model.JwtClaims).Uuid)
	if err != nil {
		zap.L().Error("GetUserByID error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "获取用户信息失败", nil))
		return nil, err
	}
	if num == 0 {
		zap.L().Error("Can not found this user")
		context.JSON(http.StatusOK, response.Result(400, "没有找到用户", nil))
		return nil, err
	}
	return userinfo, nil
}
