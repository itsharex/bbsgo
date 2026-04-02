package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("获取个人资料失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("获取用户信息失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	utils.Success(w, user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("更新个人资料失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("更新用户资料请求解析失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if password, ok := updates["password"].(string); ok && password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Printf("密码加密失败: userID=%d, 错误: %v", userID, err)
			utils.Error(w, 500, "密码加密失败")
			return
		}
		updates["password_hash"] = hashedPassword
		delete(updates, "password")
	}

	delete(updates, "id")
	delete(updates, "username")
	delete(updates, "email")
	delete(updates, "role")
	delete(updates, "credits")
	delete(updates, "level")
	delete(updates, "created_at")

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		log.Printf("更新用户资料失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("获取用户信息失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}
	utils.Success(w, user)
}

func GetUserTopics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("获取用户话题失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var topics []models.Topic
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).Preload("User").Preload("Forum").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&topics)

	utils.Success(w, map[string]interface{}{
		"list":  topics,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func GetCreditUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	database.DB.Order("credits DESC").Limit(10).Find(&users)
	utils.Success(w, users)
}
