package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"net/http"
	"strconv"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var notifications []models.Notification
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications)

	utils.Success(w, map[string]interface{}{
		"list":      notifications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetUnreadNotificationCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var count int64
	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count)

	utils.Success(w, map[string]int64{"count": count})
}

func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true)

	utils.Success(w, nil)
}

func CreateNotification(userID uint, notifType, content, link string) {
	notification := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Content: content,
		Link:    link,
		IsRead:  false,
	}
	database.DB.Create(&notification)
}
