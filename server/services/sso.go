package services

import (
	"bbsgo/database"
	"bbsgo/models"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// SSOUserInfo 第三方返回的用户信息
type SSOUserInfo struct {
	UID      int64  `json:"uid"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// VerifySSO 调用主站 verify 接口获取用户信息
// token 从浏览器 cookie 中获取，传递给主站验证
func VerifySSO(verifyURL, token string) (*SSOUserInfo, error) {
	if verifyURL == "" {
		return nil, fmt.Errorf("verify_url is empty")
	}

	// 去掉 URL 尾随空白字符
	verifyURL = strings.TrimSpace(verifyURL)

	// 构建请求，token 放到 body 中
	payload := fmt.Sprintf(`{"token":"%s"}`, token)
	req, err := http.NewRequest("POST", verifyURL, strings.NewReader(payload))
	if err != nil {
		log.Printf("[sso] create request failed: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "BBSGO-SSO/1.0")

	// 发送请求，设置超时
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[sso] request verify url failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[sso] read response failed: %v", err)
		return nil, err
	}

	// 解析 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[sso] parse response failed: %v, body: %s", err, string(body))
		return nil, fmt.Errorf("parse response failed")
	}

	// 检查错误码
	if errCode, ok := result["err_code"].(float64); ok && errCode != 200 {
		errMsg, _ := result["err_msg"].(string)
		log.Printf("[sso] verify failed: err_code=%d, err_msg=%s", int(errCode), errMsg)
		return nil, fmt.Errorf("verify failed: %s", errMsg)
	}

	// 提取用户信息（可能在顶层或 response 字段里）
	userInfo := &SSOUserInfo{}
	var data map[string]interface{}

	if respData, ok := result["response"].(map[string]interface{}); ok {
		data = respData
	} else {
		data = result
	}

	if uid, ok := data["uid"].(float64); ok {
		userInfo.UID = int64(uid)
	} else {
		return nil, fmt.Errorf("invalid uid in response")
	}

	if account, ok := data["account"].(string); ok {
		userInfo.Account = account
	}

	if nickname, ok := data["nickname"].(string); ok {
		userInfo.Nickname = nickname
	}

	if avatar, ok := data["avatar"].(string); ok {
		userInfo.Avatar = avatar
	}

	log.Printf("[sso] verify success: uid=%d, account=%s", userInfo.UID, userInfo.Account)
	return userInfo, nil
}

// GetOrCreateSSOUser 根据第三方用户信息，查找或创建本地用户
func GetOrCreateSSOUser(userInfo *SSOUserInfo) (*models.User, error) {
	var user models.User

	// 先查找是否已存在该 SSO 用户（通过 sso_uid 查找）
	err := database.DB.Where("sso_uid = ?", fmt.Sprintf("%d", userInfo.UID)).First(&user).Error
	if err == nil {
		// 找到老用户，更新信息
		log.Printf("[sso] found existing user: id=%d, uid=%s", user.ID, userInfo.UID)
		user.Nickname = userInfo.Nickname
		if userInfo.Avatar != "" {
			user.Avatar = userInfo.Avatar
		}
		if err := database.DB.Save(&user).Error; err != nil {
			log.Printf("[sso] update user failed: %v", err)
		}
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		log.Printf("[sso] query user failed: %v", err)
		return nil, err
	}

	// 创建新用户
	user = models.User{
		Username:    generateUniqueUsername(userInfo.Account),
		Email:       fmt.Sprintf("sso_%d@example.com", userInfo.UID),
		Nickname:    userInfo.Nickname,
		Avatar:      userInfo.Avatar,
		PasswordHash: generateRandomPassword(), // SSO 用户随机密码
		SSOUid:      fmt.Sprintf("%d", userInfo.UID),
		Role:        0,
		Credits:     0,
		Level:       1,
		Reputation:  100,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("[sso] create user failed: %v", err)
		return nil, err
	}

	log.Printf("[sso] created new user: id=%d, username=%s, uid=%d", user.ID, user.Username, userInfo.UID)
	return &user, nil
}

// generateUniqueUsername 生成唯一的用户名
func generateUniqueUsername(account string) string {
	// 清理用户名
	username := strings.TrimSpace(account)
	if username == "" {
		username = fmt.Sprintf("sso_%d", time.Now().UnixNano())
	}

	// 检查是否已存在
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count == 0 {
		return username
	}

	// 追加后缀直到唯一
	suffix := 1
	for {
		newUsername := fmt.Sprintf("%s_%d", username, suffix)
		database.DB.Model(&models.User{}).Where("username = ?", newUsername).Count(&count)
		if count == 0 {
			return newUsername
		}
		suffix++
		if suffix > 1000 {
			// 防止无限循环
			return fmt.Sprintf("sso_%d_%d", time.Now().UnixNano(), suffix)
		}
	}
}

// generateRandomPassword 生成随机密码
func generateRandomPassword() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "sso_" + hex.EncodeToString(bytes)
}
