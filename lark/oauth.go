package lark

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"scutbot.cn/uniauth/model"
)

// 创建 Client
// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
var client = lark.NewClient(
	viper.GetString("feishu.app_id"),
	viper.GetString("feishu.app_secret"),
	lark.WithEnableTokenCache(viper.GetBool("feishu.enableTokenCache")))

// GetUserToken 获取UserToken
func GetUserToken(code string) (string, error) {
	// 创建请求对象
	req := larkauthen.NewCreateOidcAccessTokenReqBuilder().
		Body(larkauthen.NewCreateOidcAccessTokenReqBodyBuilder().
			GrantType(`authorization_code`).
			Code(code).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Authen.OidcAccessToken.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return " ", err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		respError := errors.New(resp.Error())
		return " ", respError
	}
	// 业务处理
	// fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.Data.AccessToken)
	return *resp.Data.AccessToken, nil
}

// GetUserInfo 获取UserInfo
func GetUserInfo(userToken string) (*model.FeishuUserInfoReceiver, error) {
	resp, err := client.Authen.UserInfo.Get(context.Background(), larkcore.WithUserAccessToken(userToken))
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		respError := errors.New(resp.Error())
		return nil, respError
	}
	// 业务处理
	// fmt.Println(larkcore.Prettify(resp))
	fui := model.FeishuUserInfoReceiver(*resp.Data)
	return &fui, nil
}

// GetUserGroup 获取用户所属用户组,使用OpenID
func GetUserGroup(openId string) {
	// 构建请求
	req := larkcontact.NewMemberBelongGroupReqBuilder().
		MemberId(openId).
		MemberIdType(`open_id`).
		GroupType(1).
		PageSize(500).
		Build()
	// 发起请求
	resp, err := client.Contact.Group.MemberBelong(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}
	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}
