package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"scutbot.cn/uniauth/utils"
)

var (
	privateKeyPath = viper.Get("jwt.private_key_path").(string)
	publicKeyPath  = viper.Get("jwt.public_key_path").(string)
)

func publicKeyToJWK(publicKey *rsa.PublicKey) (*jwk.Key, error) {
	// 使用 jwk.New 创建一个新的 JSONWebKey
	jsonWebKey, err := jwk.New(publicKey)
	if err != nil {
		return nil, err
	}
	jsonWebKey.Set("kid", utils.RandStringRunes(16))
	jsonWebKey.Set("alg", "RS256")
	jsonWebKey.Set("use", "sig")
	if err != nil {
		zap.L().Error("Failed to set JWK parameters", zap.Error(err))
		return nil, err
	}
	return &jsonWebKey, nil
}

func loadKeyFromFile() (*rsa.PrivateKey, *rsa.PublicKey) {
	// 读取 PEM 文件
	privatePemBytes, err := os.ReadFile(privateKeyPath)
	publicPemBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		zap.L().Error("Failed to read key files", zap.Error(err))
		return nil, nil
	}
	// 解码 PEM 数据块
	privateBlock, _ := pem.Decode(privatePemBytes)
	publicBlock, _ := pem.Decode(publicPemBytes)
	if privateBlock == nil || privateBlock.Type != "PRIVATE KEY" {
		zap.L().Error("Failed to decode private key")
		return nil, nil
	}
	if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" {
		zap.L().Error("Failed to decode public key")
		return nil, nil
	}
	// 解析 RSA 私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		zap.L().Panic("Failed to parse private or public key")
		return nil, nil
	}

	return privateKey.(*rsa.PrivateKey), publicKey.(*rsa.PublicKey)
}
func GenerateJWK() (jwk.Key, error) {
	_, publicKey := loadKeyFromFile()

	key, err := publicKeyToJWK(publicKey)
	if err != nil {
		zap.L().Error("publicKeyToJWK error", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Generate JWK success", zap.Any("JWK", key))
	return *key, nil
}
