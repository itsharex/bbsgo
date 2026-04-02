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

type LikeRequest struct {
	TargetType string `json:"target_type"`
	TargetID   uint   `json:"target_id"`
}

func CreateLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("创建点赞失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建点赞请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.TargetType != "topic" && req.TargetType != "post" {
		log.Printf("创建点赞失败: 无效的目标类型, targetType=%s", req.TargetType)
		utils.Error(w, 400, "无效的目标类型")
		return
	}

	var existingLike models.Like
	database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).First(&existingLike)
	if existingLike.ID != 0 {
		log.Printf("创建点赞失败: 已点赞, userID=%d, targetType=%s, targetID=%d", userID, req.TargetType, req.TargetID)
		utils.Error(w, 400, "已经点赞过了")
		return
	}

	like := models.Like{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	if err := database.DB.Create(&like).Error; err != nil {
		log.Printf("创建点赞失败: userID=%d, targetType=%s, targetID=%d, 错误: %v", userID, req.TargetType, req.TargetID, err)
		utils.Error(w, 500, "点赞失败")
		return
	}

	if req.TargetType == "topic" {
		var topic models.Topic
		database.DB.First(&topic, req.TargetID)
		database.DB.Model(&topic).UpdateColumn("like_count", topic.LikeCount+1)
	} else if req.TargetType == "post" {
		var post models.Post
		database.DB.First(&post, req.TargetID)
		database.DB.Model(&post).UpdateColumn("like_count", post.LikeCount+1)
	}

	utils.Success(w, like)
}

func DeleteLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("删除点赞失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("删除点赞请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	var like models.Like
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).First(&like).Error; err != nil {
		log.Printf("查询点赞失败: userID=%d, targetType=%s, targetID=%d, 错误: %v", userID, req.TargetType, req.TargetID, err)
		utils.Error(w, 404, "点赞记录不存在")
		return
	}

	database.DB.Unscoped().Delete(&like)

	if req.TargetType == "topic" {
		var topic models.Topic
		database.DB.First(&topic, req.TargetID)
		if topic.LikeCount > 0 {
			database.DB.Model(&topic).UpdateColumn("like_count", topic.LikeCount-1)
		}
	} else if req.TargetType == "post" {
		var post models.Post
		database.DB.First(&post, req.TargetID)
		if post.LikeCount > 0 {
			database.DB.Model(&post).UpdateColumn("like_count", post.LikeCount-1)
		}
	}

	utils.Success(w, nil)
}
