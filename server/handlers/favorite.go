package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
)

type FavoriteRequest struct {
	TopicID uint `json:"topic_id"`
}

func CreateFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("创建收藏失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建收藏请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	var existingFavorite models.Favorite
	database.DB.Where("user_id = ? AND topic_id = ?", userID, req.TopicID).First(&existingFavorite)
	if existingFavorite.ID != 0 {
		log.Printf("创建收藏失败: 已收藏, userID=%d, topicID=%d", userID, req.TopicID)
		utils.Error(w, 400, "已经收藏过了")
		return
	}

	favorite := models.Favorite{
		UserID:  userID,
		TopicID: req.TopicID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		log.Printf("创建收藏失败: userID=%d, topicID=%d, 错误: %v", userID, req.TopicID, err)
		utils.Error(w, 500, "收藏失败")
		return
	}

	utils.Success(w, favorite)
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("删除收藏失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("删除收藏请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, req.TopicID).First(&favorite).Error; err != nil {
		log.Printf("查询收藏失败: userID=%d, topicID=%d, 错误: %v", userID, req.TopicID, err)
		utils.Error(w, 404, "收藏记录不存在")
		return
	}

	database.DB.Unscoped().Delete(&favorite)
	utils.Success(w, nil)
}

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("获取收藏列表失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var favorites []models.Favorite
	database.DB.Where("user_id = ?", userID).Preload("Topic").Preload("Topic.User").Preload("Topic.Forum").
		Order("created_at DESC").Find(&favorites)

	utils.Success(w, favorites)
}
