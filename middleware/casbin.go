package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"scutbot.cn/uniauth/internal"
	"scutbot.cn/uniauth/model"
)

func CheckPermission(domains ...string) gin.HandlerFunc {
	return func(context *gin.Context) {
		dbMap := viper.GetStringMapString("database")
		var parse string = dbMap["user"] + ":" + dbMap["passwd"] + "@tcp(" + dbMap["host"] + ")/" + dbMap["name"]
		a, _ := gormadapter.NewAdapter("mysql", parse, true)
		e, err := casbin.NewEnforcer("config/model.conf", a)
		if err != nil {
			zap.L().Error("Enforcer init error", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "Enforcer init error", nil))
			context.Abort()
			return
		}
		claim, _ := context.Get("claims")
		uuid := claim.(*model.JwtClaims).Uuid
		identities, err := internal.GetUserIdentityByUuID(uuid)
		if err != nil {
			zap.L().Error("Database error", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
			return
		}

		//需要组长权限的情况
		for _, domain := range domains {
			if domain == "All" {
				for _, identity := range identities {
					if identity.Role == "组长" {
						zap.L().Info("Check permission success")
						context.Next()
						return
					}
				}
				zap.L().Info("No permission")
				context.JSON(http.StatusOK, response.Result(403, "没有权限", nil))
				context.Abort()
				return
			}
		}

		//需要检查多组身份的情况
		var ok bool = false
		for _, id := range identities {
			joinTime := fmt.Sprintf("%d", id.JoinTime)
			dom := joinTime + id.Group
			for _, domain := range domains {
				if dom == domain {
					sub := id.Role                  // 想要访问资源的用户。
					obj := context.Request.URL.Path // 将被访问的资源。
					act := context.Request.Method   // 用户对资源执行的操作。
					zap.L().Info("Check permission", zap.String("sub", sub), zap.String("dom", domain), zap.String("obj", obj), zap.String("act", act))
					ok, err = e.Enforce(sub, dom, obj, act)
					if err != nil {
						zap.L().Error("Check permission failed", zap.Error(err))
						context.JSON(http.StatusOK, response.Result(500, "权限查询错误", nil))
						context.Abort()
						return
					}
					if ok {
						break
					}
				}
			}
		}

		if ok {
			zap.L().Info("Check permission success")
		} else {
			zap.L().Info("No permission")
			context.JSON(http.StatusOK, response.Result(403, "没有权限", nil))
			context.Abort()
			return
		}
	}
}
