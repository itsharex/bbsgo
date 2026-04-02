package handlers

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("注册请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Printf("注册失败: 信息不完整, username=%s, email=%s", req.Username, req.Email)
		utils.Error(w, 400, "请填写完整信息")
		return
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
		log.Printf("注册失败: 邮箱已被注册, email=%s", req.Email)
		utils.Error(w, 400, "邮箱已被注册")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("密码加密失败: 错误: %v", err)
		utils.Error(w, 500, "密码加密失败")
		return
	}

	user := models.User{
		Username:     req.Username,
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
		log.Printf("生成令牌失败: userID=%d, 错误: %v", user.ID, err)
		utils.Error(w, 500, "生成令牌失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("登录请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	var user models.User
	database.DB.Where("username = ?", req.Username).First(&user)
	if user.ID == 0 {
		log.Printf("登录失败: 用户不存在, username=%s", req.Username)
		utils.Error(w, 400, "用户名或密码错误")
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		log.Printf("登录失败: 密码错误, username=%s", req.Username)
		utils.Error(w, 400, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成令牌失败: userID=%d, 错误: %v", user.ID, err)
		utils.Error(w, 500, "生成令牌失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
