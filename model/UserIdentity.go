package model

import "github.com/jinzhu/gorm"

type UserIdentity struct {
	gorm.Model
	Uuid     string `json:"uuid"`
	Role     string `json:"role"`
	Group    string `json:"group"`
	JoinTime int    `json:"joinTime"`
	Status   int    `json:"status"` //三种状态，0代表待审批，1代表审批通过，-1代表审批不通过
}
