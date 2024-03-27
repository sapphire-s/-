package model

import "github.com/dgrijalva/jwt-go"

type JwtClaims struct {
	Uuid   string   `json:"uuid"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Groups []string `json:"group_info"`
	Avatar string   `json:"avatar"`
	jwt.StandardClaims
}
