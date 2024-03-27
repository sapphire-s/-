package model

import (
	"github.com/jinzhu/gorm"
)

type UserInfo struct {
	gorm.Model
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Avatar   string `json:"avatar"`
	Audience string `json:"aud" gorm:"type:varchar(512)"`
}
