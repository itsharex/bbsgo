package handlers

import (
	"bbsgo/config"
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"unicode"
)

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// setAuthCookie 设置认证 Cookie
func setAuthCookie(w http.ResponseWriter, token string) {
	// 获取 token 过期天数配置，默认为 7 天
	expireDays := config.GetConfigInt("jwt_expire_days", 7)
	maxAge := expireDays * 24 * 60 * 60

	// 设置 httpOnly Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "bbsgo_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   false, // 开发环境设为 false，生产环境应设为 true
		SameSite: http.SameSiteLaxMode,
	})
}

// clearAuthCookie 清除认证 Cookie
func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "bbsgo_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// Register 用户注册处理器
// 处理用户注册请求，创建新用户并设置 httpOnly Cookie
func Register(w http.ResponseWriter, r *http.Request) {
	// 检查是否允许注册
	if !config.GetConfigBool("allow_register", true) {
		log.Printf("register: registration disabled")
		errors.Error(w, errors.CodeRegisterDisabled, "")
		return
	}

	// 解析请求体
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("register: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Printf("register: incomplete registration info, username: %s, email: %s", req.Username, req.Email)
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	// 验证用户名格式
	if err := utils.ValidateUsername(req.Username); err != nil {
		log.Printf("register: invalid username, username: %s, error: %v", req.Username, err)
		errors.Error(w, errors.CodeInvalidUsername, err.Error())
		return
	}

	// 验证邮箱格式
	if err := utils.ValidateEmail(req.Email); err != nil {
		log.Printf("register: invalid email, email: %s, error: %v", req.Email, err)
		errors.Error(w, errors.CodeInvalidEmail, err.Error())
		return
	}

	// 验证密码强度：长度至少8位，包含至少3种字符类型
	if err := validatePasswordSimple(req.Password); err != nil {
		log.Printf("register: password validation failed, username: %s, error: %v", req.Username, err)
		errors.Error(w, errors.CodePasswordTooWeak, err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		log.Printf("register: username already exists, username: %s", req.Username)
		errors.Error(w, errors.CodeUsernameExists, "")
		return
	}

	// 检查邮箱是否已被注册
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		log.Printf("register: email already registered, email: %s", req.Email)
		errors.Error(w, errors.CodeEmailExists, "")
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("register: failed to hash password, error: %v", err)
		errors.Error(w, errors.CodePasswordHashFailed, "")
		return
	}

	// 创建用户
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         0, // 普通用户
		Credits:      0, // 初始积分
		Level:        1, // 初始等级
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("register: failed to create user, username: %s, email: %s, error: %v", req.Username, req.Email, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.TokenVersion)
	if err != nil {
		log.Printf("register: failed to generate token, userID: %d, error: %v", user.ID, err)
		errors.Error(w, errors.CodeTokenGenerateFailed, "")
		return
	}

	// 设置 httpOnly Cookie
	log.Printf("register: setting cookie for userID: %d, token: %s...", user.ID, token[:20])
	setAuthCookie(w, token)

	log.Printf("register: user registered successfully, userID: %d, username: %s", user.ID, req.Username)
	errors.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Login 用户登录处理器
// 验证用户名密码，设置 httpOnly Cookie
func Login(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("login: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 检查登录尝试限制
	allowed, remainingTime := utils.CheckLoginAttempt(req.Username)
	if !allowed {
		log.Printf("login: account temporarily locked, username: %s, remaining: %v", req.Username, remainingTime)
		errors.Error(w, errors.CodeTooManyRequests, fmt.Sprintf("账户已被临时锁定，请 %v 后再试", remainingTime.Round(time.Minute)))
		return
	}

	// 查询用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("login: user not found, username: %s", req.Username)
		// 记录登录失败
		utils.RecordLoginFailure(req.Username)
		errors.Error(w, errors.CodeUsernameOrPassword, "")
		return
	}

	// 检查用户是否被封禁
	if user.IsBanned {
		log.Printf("login: user is banned, username: %s", req.Username)
		errors.Error(w, errors.CodeUserBanned, "您的账户已被封禁")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		log.Printf("login: password mismatch, username: %s", req.Username)
		// 记录登录失败
		utils.RecordLoginFailure(req.Username)

		// 获取剩余尝试次数
		attempts := utils.GetLoginAttempts(req.Username)
		remaining := 5 - attempts
		if remaining < 0 {
			remaining = 0
		}

		if remaining > 0 {
			errors.Error(w, errors.CodeUsernameOrPassword, fmt.Sprintf("用户名或密码错误，剩余尝试次数: %d", remaining))
		} else {
			errors.Error(w, errors.CodeUsernameOrPassword, "")
		}
		return
	}

	// 登录成功，清除失败记录
	utils.RecordLoginSuccess(req.Username)

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.TokenVersion)
	if err != nil {
		log.Printf("login: failed to generate token, userID: %d, error: %v", user.ID, err)
		errors.Error(w, errors.CodeTokenGenerateFailed, "")
		return
	}

	// 设置 httpOnly Cookie
	setAuthCookie(w, token)

	log.Printf("login: user logged in successfully, userID: %d, username: %s", user.ID, req.Username)
	errors.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Logout 用户登出处理器
// 清除认证 Cookie
func Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookie(w)
	log.Printf("logout: user logged out successfully")
	errors.Success(w, map[string]string{
		"message": "登出成功",
	})
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 从 context 获取用户ID
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.Error(w, errors.CodeUnauthorized, "")
		return
	}

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get current user: user not found, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	errors.Success(w, user)
}

// validatePasswordSimple 验证密码强度：长度至少8位，包含至少3种字符类型
func validatePasswordSimple(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少8位")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	if count < 3 {
		return fmt.Errorf("密码需包含大写字母、小写字母、数字和特殊字符中的至少3种")
	}

	return nil
}
