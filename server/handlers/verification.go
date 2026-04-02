package handlers

import (
	"bbsgo/config"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/services"
	"bbsgo/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type SendCodeRequest struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type RegisterWithCodeRequest struct {
	Username        string `json:"username"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Code            string `json:"code"`
}

func SendVerificationCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("发送验证码请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Email == "" {
		log.Printf("发送验证码失败: 邮箱为空")
		utils.Error(w, 400, "邮箱不能为空")
		return
	}

	emailEnabled := config.GetConfigBool("email_enabled", false)
	if !emailEnabled {
		log.Printf("发送验证码失败: 邮件服务未启用")
		utils.Error(w, 500, "验证码发送失败，请重试")
		return
	}

	var existingUser models.User
	database.DB.Where("email = ?", req.Email).First(&existingUser)
	if existingUser.ID != 0 {
		log.Printf("发送验证码失败: 邮箱已被注册, email=%s", req.Email)
		utils.Error(w, 400, "该邮箱已被注册")
		return
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(5 * time.Minute)

	verificationCode := models.VerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      "register",
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&verificationCode).Error; err != nil {
		log.Printf("保存验证码失败: email=%s, 错误: %v", req.Email, err)
		utils.Error(w, 500, "发送验证码失败")
		return
	}

	if err := services.SendVerificationCode(req.Email, code); err != nil {
		log.Printf("发送邮件失败: email=%s, 错误: %v", req.Email, err)
		utils.Error(w, 500, "发送验证码失败")
		return
	}

	log.Printf("发送验证码成功: email=%s", req.Email)
	utils.Success(w, map[string]string{
		"message": "验证码已发送",
	})
}

func RegisterWithCode(w http.ResponseWriter, r *http.Request) {
	var req RegisterWithCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("注册请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Username == "" || req.Nickname == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		log.Printf("注册失败: 信息不完整, username=%s, nickname=%s, email=%s", req.Username, req.Nickname, req.Email)
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	if req.Password != req.ConfirmPassword {
		log.Printf("注册失败: 两次密码不一致, username=%s", req.Username)
		utils.Error(w, 400, "两次密码输入不一致")
		return
	}

	emailEnabled := config.GetConfigBool("email_enabled", false)
	if emailEnabled {
		if req.Code == "" {
			log.Printf("注册失败: 验证码为空, username=%s", req.Username)
			utils.Error(w, 400, "请输入验证码")
			return
		}

		var verificationCode models.VerificationCode
		result := database.DB.Where("email = ? AND code = ? AND type = ? AND expires_at > ?",
			req.Email, req.Code, "register", time.Now()).First(&verificationCode)
		if result.Error != nil {
			log.Printf("注册失败: 验证码无效或已过期, username=%s, email=%s, code=%s", req.Username, req.Email, req.Code)
			utils.Error(w, 400, "验证码无效或已过期")
			return
		}

		database.DB.Unscoped().Delete(&verificationCode)
	}

	var existingUser models.User
	database.DB.Where("username = ?", req.Username).First(&existingUser)
	if existingUser.ID != 0 {
		log.Printf("注册失败: 用户名已存在, username=%s", req.Username)
		utils.Error(w, 400, "用户名已存在")
		return
	}

	database.DB.Where("email = ?", req.Email).First(&existingUser)
	if existingUser.ID != 0 {
		log.Printf("注册失败: 邮箱已被注册, username=%s, email=%s", req.Username, req.Email)
		utils.Error(w, 400, "邮箱已被注册")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("密码加密失败: username=%s, 错误: %v", req.Username, err)
		utils.Error(w, 500, "密码加密失败")
		return
	}

	user := models.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         0,
		Credits:      0,
		Level:        1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("创建用户失败: username=%s, email=%s, 错误: %v", req.Username, req.Email, err)
		utils.Error(w, 500, "注册失败")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成令牌失败: userID=%d, username=%s, 错误: %v", user.ID, user.Username, err)
		utils.Error(w, 500, "生成令牌失败")
		return
	}

	log.Printf("注册成功: userID=%d, username=%s, email=%s", user.ID, user.Username, user.Email)
	utils.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
