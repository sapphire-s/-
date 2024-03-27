package model

type FeishuUserInfo struct {
	UnionId   string `json:"union_id" gorm:"primaryKey;not null;"` // 用户统一ID
	UserId    string `json:"user_id" gorm:"unique;not null;"`      // 用户 user_id
	OpenId    string `json:"open_id" gorm:"unique;not null;"`      // 用户在应用内的唯一标识
	Name      string `json:"name"`                                 // 用户姓名
	Email     string `json:"email" gorm:"unique"`                  // 用户邮箱
	Mobile    string `json:"mobile" gorm:"unique"`                 // 用户手机号
	AvatarUrl string `json:"avatar_url"`                           // 用户头像
	AvatarBig string `json:"avatar_big"`                           // 用户头像 640x640
}
