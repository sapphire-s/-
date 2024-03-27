package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"scutbot.cn/uniauth/internal"
	"scutbot.cn/uniauth/model"
	"strconv"
)

func Apply(context *gin.Context) {
	identity := model.UserIdentity{}
	err := context.ShouldBindJSON(&identity)
	if err != nil {
		zap.L().Error("Parse info error", zap.Error(err), zap.Any("identity", identity))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	if !isCorrect(identity) {
		zap.L().Error("Invalid data", zap.Error(err), zap.Any("identity", identity))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	claim, _ := context.Get("claims")
	uuid := claim.(*model.JwtClaims).Uuid
	result, err := internal.GetUserIdentityByUuID(uuid)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	for _, userIdentity := range result {
		if userIdentity.Group == identity.Group && userIdentity.Role == identity.Role && userIdentity.JoinTime == identity.JoinTime {
			zap.L().Info("Identity exist", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(400, "已存在该身份", err))
			return
		}
	}
	identity.Status = 0
	identity.Uuid = uuid
	err = internal.AddUserIdentity(&identity)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	context.JSON(http.StatusOK, response.Result(200, "身份申请成功", nil))
}
func Approve(context *gin.Context) {
	identity := model.UserIdentity{}
	err := context.ShouldBindJSON(&identity)
	if err != nil {
		zap.L().Error("Parse info error", zap.Error(err), zap.Any("identity", identity))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	status := context.Query("status")
	statusInt, err := strconv.ParseInt(status, 10, 8)
	if err != nil {
		zap.L().Error("Parse status error", zap.Error(err), zap.Any("status", status))
		context.JSON(http.StatusOK, response.Result(500, "status error", err))
		return
	}
	id, err := internal.GetUserIdentityByID(identity.ID)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	id.Status = int(statusInt)
	err = internal.UpdateIdentity(id)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	if statusInt == 1 {
		context.JSON(http.StatusOK, response.Result(200, "审批已同意", nil))
	} else {
		context.JSON(http.StatusOK, response.Result(200, "审批已拒绝", nil))
	}
}
func Withdraw(context *gin.Context) {
	identity := model.UserIdentity{}
	err := context.ShouldBindJSON(&identity)
	if err != nil {
		zap.L().Error("Parse info error", zap.Error(err), zap.Any("identity", identity))
		context.JSON(http.StatusOK, response.Result(400, "错误请求", err))
		return
	}
	err = internal.DeleteUserIdentity(identity.ID)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	context.JSON(http.StatusOK, response.Result(200, "审批已撤回/删除", nil))
}

func GetIdentity(context *gin.Context) {
	claim, _ := context.Get("claims")
	uuid := claim.(*model.JwtClaims).Uuid
	preResult, err := internal.GetUserIdentityByUuID(uuid)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	var result []model.UserIdentity
	for _, r := range preResult {
		result = append(result, *r)
	}
	context.JSON(http.StatusOK, response.Result(200, "获取成功", result))
}
func Pending(context *gin.Context) {
	//获取请求的用户身份
	claim, _ := context.Get("claims")
	uuid := claim.(*model.JwtClaims).Uuid
	userIdentities, err := internal.GetUserIdentityByUuID(uuid)
	if err != nil {
		zap.L().Error("Database error", zap.Error(err))
		context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
		return
	}
	var pendingGroups []string
	for _, identity := range userIdentities {
		if identity.Status == 1 {
			if identity.Role == "组长" {
				pendingGroups = append(pendingGroups, identity.Group)
			}
		}
	}
	zap.L().Info("pendingGroups", zap.Strings("pendingGroups", pendingGroups))
	//根据用户拥有的组长身份，寻找相应的待审批的身份
	var pendingIdentities []*model.UserIdentity
	for _, group := range pendingGroups {
		zap.L().Info("group", zap.String("group", group))
		temp, err := internal.GetPendingUserIdentityByGroup(group)
		if err != nil {
			zap.L().Error("Database error", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
			return
		}
		for _, identity := range temp {
			pendingIdentities = append(pendingIdentities, identity)
		}
	}
	//把Uuid变成用户名
	for _, identity := range pendingIdentities {
		userinfo, num, err := internal.GetUserByID(identity.Uuid)
		if err != nil {
			zap.L().Error("Database error", zap.Error(err))
			context.JSON(http.StatusOK, response.Result(500, "数据库错误", err))
			return
		}
		if num == 0 {
			zap.L().Error("Uuid error", zap.Any("userinfo", userinfo))
			context.JSON(http.StatusOK, response.Result(400, "用户不存在", nil))
			return
		}
		identity.Uuid = userinfo.Name
	}
	context.JSON(http.StatusOK, response.Result(200, "获取成功", pendingIdentities))
}
func isCorrect(identity model.UserIdentity) bool {
	var legalGroups []string = []string{"机械组", "电控组", "视觉组", "软开组", "宣运组", "管理组", "教育组", "项目组", "顾问组"}
	var legalRoles []string = []string{"组长", "组员"}
	group := false
	role := false
	for _, legalGroup := range legalGroups {
		if legalGroup == identity.Group {
			group = true
			break
		}
	}
	for _, legalRole := range legalRoles {
		if legalRole == identity.Role {
			role = true
			break
		}
	}
	if group && role {
		return true
	}
	return false
}
