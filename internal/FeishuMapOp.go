package internal

import "scutbot.cn/uniauth/model"

func AddFeishuMap(feishuMap *model.FeishuMap) error {
	db.AutoMigrate(feishuMap)
	return db.Create(feishuMap).Error
}
func GetFeishuMapByUnionId(unionId string) (*model.FeishuMap, int64, error) {
	var feishuMap model.FeishuMap
	result := db.Where("union_id = ?", unionId).First(&feishuMap)
	return &feishuMap, result.RowsAffected, result.Error
}
func DeleteFeishuMap(feishuMap *model.FeishuMap) error {
	return db.Delete(feishuMap).Error
}

func GetUuid(union_id string) (string, error) {
	var feishuMap model.FeishuMap
	result := db.Where("union_id = ?", union_id).First(&feishuMap)
	return feishuMap.Uuid, result.Error
}
