package internal

import "scutbot.cn/uniauth/model"

// NewFeishuUserInfo 通过reciver创建FeishuUserInfo
func NewFeishuUserInfo(receiver *model.FeishuUserInfoReceiver) *model.FeishuUserInfo {
	result := model.FeishuUserInfo{
		UnionId:   *receiver.UnionId,
		UserId:    *receiver.UserId,
		OpenId:    *receiver.OpenId,
		Name:      *receiver.Name,
		Email:     *receiver.Email,
		Mobile:    *receiver.Mobile,
		AvatarUrl: *receiver.AvatarUrl,
		AvatarBig: *receiver.AvatarBig,
	}
	return &result
}

// AddFeishuUserInfo 添加FeishuUserInfo
func AddFeishuUserInfo(feishuUserInfo *model.FeishuUserInfo) error {
	return db.Create(feishuUserInfo).Error
}

// GetFeishuUserInfoByUnionId 通过unionId获取FeishuUserInfo
func GetFeishuUserInfoByUnionId(unionId string) (*model.FeishuUserInfo, int64, error) {
	var feishuUserInfo model.FeishuUserInfo
	result := db.Where("union_id = ?", unionId).First(&feishuUserInfo)
	return &feishuUserInfo, result.RowsAffected, result.Error
}

// GetFeishuUserInfoByOpenId 通过openId获取FeishuUserInfo
func GetFeishuUserInfoByOpenId(openId string) (*model.FeishuUserInfo, error) {
	var feishuUserInfo model.FeishuUserInfo
	return &feishuUserInfo, db.Where("open_id = ?", openId).First(&feishuUserInfo).Error
}

// UpdateFeishuUserInfo 更新FeishuUserInfo
func UpdateFeishuUserInfo(feishuUserInfo *model.FeishuUserInfo) error {
	return db.Save(feishuUserInfo).Error
}
