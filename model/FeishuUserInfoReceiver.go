package model

type FeishuUserInfoReceiver struct {
	Name            *string `json:"name,omitempty"`             // 用户姓名
	EnName          *string `json:"en_name,omitempty"`          // 用户英文名称
	AvatarUrl       *string `json:"avatar_url,omitempty"`       // 用户头像
	AvatarThumb     *string `json:"avatar_thumb,omitempty"`     // 用户头像 72x72
	AvatarMiddle    *string `json:"avatar_middle,omitempty"`    // 用户头像 240x240
	AvatarBig       *string `json:"avatar_big,omitempty"`       // 用户头像 640x640
	OpenId          *string `json:"open_id,omitempty"`          // 用户在应用内的唯一标识
	UnionId         *string `json:"union_id,omitempty"`         // 用户统一ID
	Email           *string `json:"email,omitempty"`            // 用户邮箱
	EnterpriseEmail *string `json:"enterprise_email,omitempty"` // 企业邮箱，请先确保已在管理后台启用飞书邮箱服务
	UserId          *string `json:"user_id,omitempty"`          // 用户 user_id
	Mobile          *string `json:"mobile,omitempty"`           // 用户手机号
	TenantKey       *string `json:"tenant_key,omitempty"`       // 当前企业标识
	EmployeeNo      *string `json:"employee_no,omitempty"`      // 用户工号
}
