package model

type FeishuMap struct {
	Uuid    string `json:"uuid" gorm:"primaryKey;not null;"`
	UnionId string `json:"unionId" gorm:"unique;not null;"`
}
