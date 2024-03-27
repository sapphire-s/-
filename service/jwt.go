package service

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"scutbot.cn/uniauth/internal"
	"scutbot.cn/uniauth/model"
)

// JWT jwt签名结构
type JWT struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

var (
	priKey, pubKey = loadKeyFromFile()
	issuer         = viper.Get("jwt.issuer").(string)
)

func NewJWT() *JWT {
	return &JWT{
		publicKey:  pubKey,
		privateKey: priKey,
	}
}

// 自定义错误类型
var (
	TokenInvalid     = errors.New("token is invalid")
	TokenExpired     = errors.New("token is expired")
	TokenMalformed   = errors.New("token is malformed")
	TokenNotValidYet = errors.New("token is not valid yet")
)

// CreateToken 创建Token
func (j *JWT) CreateToken(claims model.JwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*model.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				fmt.Println(err.Error())
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*model.JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// RefreshToken 更新Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &model.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*model.JwtClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// GenerateToken 生成token
func GenerateToken(usr *model.UserInfo, aud string) (string, error) {
	userIdentity, err := internal.GetUserIdentityByUuID(usr.Uuid)
	if err != nil {
		return "", err
	}
	var Groups []string
	for _, identity := range userIdentity {
		Groups = append(Groups, identity.Group)
	}

	j := NewJWT()
	claims := model.JwtClaims{
		Uuid:   usr.Uuid,
		Email:  usr.Email,
		Groups: Groups,
		Name:   usr.Name,
		Avatar: usr.Avatar,
		StandardClaims: jwt.StandardClaims{
			Audience:  aud,                              // client_id
			NotBefore: int64(time.Now().Unix() - 1000),  // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 21600), // 签名过期时间 六小时
			IssuedAt:  int64(time.Now().Unix()),         // 签名生成时间
			Issuer:    issuer,                           // 签名发行者
			Subject:   usr.Uuid,                         // 用户标识符
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
