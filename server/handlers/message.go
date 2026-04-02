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

func GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var messages []models.Message
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.Message{}).Where("to_user_id = ? OR from_user_id = ?", userID, userID).Count(&total)
	database.DB.Where("to_user_id = ? OR from_user_id = ?", userID, userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Preload("FromUser").Preload("ToUser").
		Find(&messages)

	utils.Success(w, map[string]interface{}{
		"list":      messages,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetMessageConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	otherUserID, _ := strconv.Atoi(vars["user_id"])

	var messages []models.Message
	database.DB.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userID, otherUserID, otherUserID, userID).
		Order("created_at ASC").
		Preload("FromUser").Preload("ToUser").
		Find(&messages)

	database.DB.Model(&models.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND is_read = ?", otherUserID, userID, false).
		Update("is_read", true)

	utils.Success(w, messages)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var req struct {
		ToUserID uint   `json:"to_user_id"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("发送消息请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Content == "" {
		log.Printf("发送消息失败: 内容为空, fromUserID=%d, toUserID=%d", userID, req.ToUserID)
		utils.Error(w, 400, "消息内容不能为空")
		return
	}

	var toUser models.User
	if err := database.DB.First(&toUser, req.ToUserID).Error; err != nil {
		log.Printf("获取接收用户失败: toUserID=%d, 错误: %v", req.ToUserID, err)
		utils.Error(w, 404, "接收用户不存在")
		return
	}

	message := models.Message{
		FromUserID: userID,
		ToUserID:   req.ToUserID,
		Content:    req.Content,
		IsRead:     false,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		log.Printf("创建消息失败: fromUserID=%d, toUserID=%d, 错误: %v", userID, req.ToUserID, err)
		utils.Error(w, 500, "发送失败")
		return
	}

	notification := models.Notification{
		UserID:  req.ToUserID,
		Type:    "message",
		Content: "您收到了一条新私信",
		Link:    "/messages",
		IsRead:  false,
	}
	database.DB.Create(&notification)

	utils.Success(w, message)
}

func MarkMessagesRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var req struct {
		FromUserID uint `json:"from_user_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	database.DB.Model(&models.Message{}).
		Where("to_user_id = ? AND from_user_id = ? AND is_read = ?", userID, req.FromUserID, false).
		Update("is_read", true)

	utils.Success(w, nil)
}

func GetUnreadMessageCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var count int64
	database.DB.Model(&models.Message{}).Where("to_user_id = ? AND is_read = ?", userID, false).Count(&count)

	utils.Success(w, map[string]int64{"count": count})
}
