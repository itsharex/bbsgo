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
	"time"

	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var posts []models.Post
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.Post{}).Where("topic_id = ? AND parent_id IS NULL", topicID).Count(&total)
	database.DB.Where("topic_id = ? AND parent_id IS NULL", topicID).Preload("User").
		Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&posts)

	utils.Success(w, map[string]interface{}{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	var req struct {
		Content  string `json:"content"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建评论请求解析失败: topicID=%d, 错误: %v", topicID, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Content == "" {
		log.Printf("创建评论失败: 内容为空, topicID=%d", topicID)
		utils.Error(w, 400, "请填写内容")
		return
	}

	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		log.Printf("获取话题失败: topicID=%d, 错误: %v", topicID, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	if topic.IsLocked || !topic.AllowComment {
		log.Printf("创建评论失败: 话题已关闭评论, topicID=%d", topicID)
		utils.Error(w, 400, "该话题已关闭评论")
		return
	}

	post := models.Post{
		TopicID:  uint(topicID),
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		log.Printf("创建评论失败: topicID=%d, userID=%d, 错误: %v", topicID, userID, err)
		utils.Error(w, 500, "发布失败")
		return
	}

	now := time.Now()
	database.DB.Model(&topic).Updates(map[string]interface{}{
		"reply_count":   topic.ReplyCount + 1,
		"last_reply_at": now,
	})

	database.DB.Preload("User").First(&post, post.ID)
	utils.Success(w, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		log.Printf("获取评论失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "评论不存在")
		return
	}

	if post.UserID != userID {
		log.Printf("编辑评论失败: 无权限, postID=%d, userID=%d", id, userID)
		utils.Error(w, 403, "无权限编辑")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("更新评论请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "topic_id")
	delete(updates, "parent_id")
	delete(updates, "created_at")
	delete(updates, "like_count")

	if err := database.DB.Model(&post).Updates(updates).Error; err != nil {
		log.Printf("更新评论失败: id=%d, userID=%d, 错误: %v", id, userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	database.DB.Preload("User").First(&post, id)
	utils.Success(w, post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		log.Printf("获取评论失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "评论不存在")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("获取用户失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	if post.UserID != userID && user.Role < 1 {
		utils.Error(w, 403, "无权限删除")
		return
	}

	var topic models.Topic
	database.DB.First(&topic, post.TopicID)
	database.DB.Unscoped().Delete(&post)
	database.DB.Model(&topic).UpdateColumn("reply_count", topic.ReplyCount-1)

	utils.Success(w, nil)
}
