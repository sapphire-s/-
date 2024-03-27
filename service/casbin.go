package service

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func AddPermission(context *gin.Context) {
	dbMap := viper.GetStringMapString("database")
	var parse string = dbMap["user"] + ":" + dbMap["passwd"] + "@tcp(" + dbMap["host"] + ")/" + dbMap["name"]
	a, _ := gormadapter.NewAdapter("mysql", parse, true)
	e, _ := casbin.NewEnforcer("config/model.conf", a)
	addPolicy(e, "组长", "24软开组", "/casbin/test", "POST", "allow")
	addPolicy(e, "组长", "24软开组", "/casbin/test", "GET", "deny")
	context.JSON(200, "添加成功")
}

func addPolicy(e *casbin.Enforcer, sub string, dom string, obj string, act string, eft string) {
	b, err := e.AddPolicy(sub, dom, obj, act, eft)
	policy := "(" + sub + "," + dom + "," + obj + "," + act + "," + eft + ")"
	if err != nil {
		zap.L().Error("添加策略失败", zap.String("policy", policy))
		return
	}
	if !b {
		zap.L().Info("已添加的策略", zap.String("policy", policy))
	} else {
		zap.L().Info("添加策略成功", zap.String("policy", policy))
	}

}
