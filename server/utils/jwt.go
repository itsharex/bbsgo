package utils

import (
	"bbsgo/config"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string
var secretOnce sync.Once

// initJWTSecret 懒初始化 JWT 密钥（在第一次使用时调用）
func initJWTSecret() {
	secretOnce.Do(func() {
		// 优先尝试从配置文件读取（不依赖数据库）
		secretFile := getSecretFilePath()
		if data, err := ioutil.ReadFile(secretFile); err == nil {
			secret := string(data)
			if len(secret) >= 32 {
				jwtSecret = secret
				log.Printf("jwt: using secret from file (length: %d)", len(secret))
				return
			}
		}

		// 生成新的随机密钥
		secret := generateRandomSecret()
		jwtSecret = secret

		// 保存到文件
		if err := saveSecretToFile(secretFile, secret); err != nil {
			log.Printf("jwt: warning: failed to save secret to file: %v", err)
		} else {
			log.Printf("jwt: generated and saved new secret to file (length: %d)", len(secret))
		}
	})
}

// getSecretFilePath 获取密钥文件路径
func getSecretFilePath() string {
	// 尝试多个可能的路径
	paths := []string{
		"./jwt.secret",
		"./data/jwt.secret",
		"./config/jwt.secret",
	}

	for _, path := range paths {
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); err == nil {
			return path
		}
	}

	// 默认使用当前目录
	return "./jwt.secret"
}

// generateRandomSecret 生成随机密钥
func generateRandomSecret() string {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("jwt: failed to generate random secret: %v", err)
	}
	return hex.EncodeToString(bytes)
}

// saveSecretToFile 保存密钥到文件
func saveSecretToFile(path, secret string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// 写入文件，权限设置为只有所有者可读写
	return ioutil.WriteFile(path, []byte(secret), 0600)
}

// Claims JWT 载荷结构
// 包含用户的基本身份信息
type Claims struct {
	UserID               uint   `json:"user_id"`       // 用户ID
	Username             string `json:"username"`      // 用户名
	TokenVersion         int    `json:"token_version"` // Token版本号，用于密码修改后使旧token失效
	jwt.RegisteredClaims        // JWT 标准声明（过期时间、签发时间等）
}

// GenerateToken 生成 JWT 令牌
// userID: 用户ID
// username: 用户名
// tokenVersion: Token版本号
// 返回: 生成的令牌字符串和错误信息
func GenerateToken(userID uint, username string, tokenVersion int) (string, error) {
	// 确保密钥已初始化
	initJWTSecret()

	// 获取 token 过期天数配置，默认为 7 天
	expireDays := config.GetConfigInt("jwt_expire_days", 7)

	// 构建 JWT Claims
	claims := Claims{
		UserID:       userID,
		Username:     username,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expireDays)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                           // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                           // 生效时间
		},
	}

	// 使用 HS256 算法签名生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken 解析并验证 JWT 令牌
// tokenString: 令牌字符串
// 返回: 解析后的 Claims 和错误信息
func ParseToken(tokenString string) (*Claims, error) {
	// 确保密钥已初始化
	initJWTSecret()

	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	// 解析失败
	if err != nil {
		log.Printf("parse token failed, error: %v", err)
		return nil, err
	}

	// 验证 token 有效性并提取 claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	// token 无效
	log.Printf("token is invalid")
	return nil, errors.New("invalid token")
}

// GetSecret 获取当前使用的密钥（仅用于调试）
func GetSecret() string {
	initJWTSecret()
	if len(jwtSecret) > 10 {
		return jwtSecret[:10] + "..." + jwtSecret[len(jwtSecret)-10:]
	}
	return "***"
}
