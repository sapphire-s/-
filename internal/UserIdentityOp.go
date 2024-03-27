package internal

import (
	"scutbot.cn/uniauth/model"
)

// NewUserIdentity 创建用户身份
func NewUserIdentity(uuid string, role string, group string, joinTime int) *model.UserIdentity {
	userIdentity := &model.UserIdentity{
		Uuid:     uuid,
		Role:     role,
		Group:    group,
		JoinTime: joinTime,
	}
	return userIdentity
}

// AddUserIdentity 添加用户身份
func AddUserIdentity(userIdentity *model.UserIdentity) error {
	// db.AutoMigrate(&User{})
	return db.Create(userIdentity).Error
}

// DeleteUserIdentity 删除用户身份
func DeleteUserIdentity(id uint) error {
	return db.Where("id = ?", id).Delete(&model.UserIdentity{}).Error
}

// 按照指定的uuid查询用户有效身份
func GetUserIdentityByUuID(uuid string) ([]*model.UserIdentity, error) {
	var UserIdentity []*model.UserIdentity
	return UserIdentity, db.Where("uuid = ?", uuid).Where("status = ?", 1).Find(&UserIdentity).Error
}

// 按照指定的id查询用户身份
func GetUserIdentityByID(id uint) (*model.UserIdentity, error) {
	var UserIdentity model.UserIdentity
	return &UserIdentity, db.Where("id = ?", id).First(&UserIdentity).Error
}

// 按照组别查询用户有效身份
func GetUserIdentityByGroup(group string) ([]model.UserIdentity, error) {
	var UserIdentity []model.UserIdentity
	return UserIdentity, db.Where("group = ?", group).Where("status = ?", 1).Find(&UserIdentity).Error
}

// 按照组别查询用户待审核身份
func GetPendingUserIdentityByGroup(group string) ([]*model.UserIdentity, error) {
	var UserIdentity []*model.UserIdentity
	return UserIdentity, db.Where("`group` = ? AND status = ?", group, 0).Find(&UserIdentity).Error
}
func UpdateIdentity(identity *model.UserIdentity) error {
	return db.Save(identity).Error
}
