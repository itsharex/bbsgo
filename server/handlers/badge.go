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

	"github.com/gorilla/mux"
)

func GetBadges(w http.ResponseWriter, r *http.Request) {
	var badges []models.Badge
	database.DB.Find(&badges)

	utils.Success(w, badges)
}

func GetUserBadges(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var userBadges []models.UserBadge
	database.DB.Where("user_id = ?", userID).Preload("Badge").Find(&userBadges)

	utils.Success(w, userBadges)
}

func CreateBadge(w http.ResponseWriter, r *http.Request) {
	var badge models.Badge
	if err := json.NewDecoder(r.Body).Decode(&badge); err != nil {
		log.Printf("创建勋章请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if badge.Name == "" {
		log.Printf("创建勋章失败: 勋章名称为空")
		utils.Error(w, 400, "请填写勋章名称")
		return
	}

	if err := database.DB.Create(&badge).Error; err != nil {
		log.Printf("创建勋章失败: name=%s, 错误: %v", badge.Name, err)
		utils.Error(w, 500, "创建失败")
		return
	}

	utils.Success(w, badge)
}

func AwardBadge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  uint `json:"user_id"`
		BadgeID uint `json:"badge_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("授予勋章请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		log.Printf("获取用户失败: userID=%d, 错误: %v", req.UserID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	var badge models.Badge
	if err := database.DB.First(&badge, req.BadgeID).Error; err != nil {
		log.Printf("获取勋章失败: badgeID=%d, 错误: %v", req.BadgeID, err)
		utils.Error(w, 404, "勋章不存在")
		return
	}

	var existing models.UserBadge
	if err := database.DB.Where("user_id = ? AND badge_id = ?", req.UserID, req.BadgeID).First(&existing).Error; err == nil {
		log.Printf("授予勋章失败: 用户已获得该勋章, userID=%d, badgeID=%d", req.UserID, req.BadgeID)
		utils.Error(w, 400, "用户已获得该勋章")
		return
	}

	userBadge := models.UserBadge{
		UserID:  req.UserID,
		BadgeID: req.BadgeID,
	}

	if err := database.DB.Create(&userBadge).Error; err != nil {
		log.Printf("授予用户勋章失败: userID=%d, badgeID=%d, 错误: %v", req.UserID, req.BadgeID, err)
		utils.Error(w, 500, "授予失败")
		return
	}

	utils.Success(w, userBadge)
}

func DeleteBadge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var badge models.Badge
	if err := database.DB.First(&badge, id).Error; err != nil {
		log.Printf("获取勋章失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "勋章不存在")
		return
	}

	database.DB.Unscoped().Delete(&badge)
	utils.Success(w, nil)
}
