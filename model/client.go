package model

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Name         string `json:"name"`
}
