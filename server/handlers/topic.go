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

func GetTopics(w http.ResponseWriter, r *http.Request) {
	forumID, _ := strconv.Atoi(r.URL.Query().Get("forum_id"))
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	sort := r.URL.Query().Get("sort")
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

	query := database.DB.Model(&models.Topic{})
	if forumID > 0 {
		query = query.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	dbQuery := database.DB.Preload("User").Preload("Forum").Preload("Tags")
	if forumID > 0 {
		dbQuery = dbQuery.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		dbQuery = dbQuery.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	switch sort {
	case "hot":
		dbQuery = dbQuery.Order("is_pinned DESC, (like_count + reply_count * 2) DESC, created_at DESC")
	case "reply":
		dbQuery = dbQuery.Order("is_pinned DESC, last_reply_at DESC NULLS LAST, created_at DESC")
	default:
		dbQuery = dbQuery.Order("is_pinned DESC, created_at DESC")
	}

	dbQuery.Offset(offset).Limit(pageSize).Find(&topics)

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var topic models.Topic
	if err := database.DB.Preload("User").Preload("Forum").First(&topic, id).Error; err != nil {
		log.Printf("获取话题失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	database.DB.Model(&topic).UpdateColumn("view_count", topic.ViewCount+1)

	utils.Success(w, topic)
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("创建话题失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	var req struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		ForumID  uint     `json:"forum_id"`
		TagNames []string `json:"tag_names"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建话题请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Title == "" || req.Content == "" || req.ForumID == 0 {
		log.Printf("创建话题失败: 信息不完整, title=%s, forumID=%d", req.Title, req.ForumID)
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	if len(req.TagNames) > 3 {
		log.Printf("创建话题失败: 标签数量超过限制, tagCount=%d", len(req.TagNames))
		utils.Error(w, 400, "最多只能添加3个标签")
		return
	}

	topic := models.Topic{
		Title:        req.Title,
		Content:      req.Content,
		UserID:       userID,
		ForumID:      req.ForumID,
		AllowComment: true,
	}

	if err := database.DB.Create(&topic).Error; err != nil {
		log.Printf("创建话题失败: userID=%d, forumID=%d, 错误: %v", userID, req.ForumID, err)
		utils.Error(w, 500, "发布失败")
		return
	}

	if len(req.TagNames) > 0 {
		var tags []models.Tag
		for _, name := range req.TagNames {
			tag, err := GetOrCreateTagByName(name)
			if err != nil || tag == nil {
				continue
			}
			if tag.IsBanned {
				continue
			}
			tags = append(tags, *tag)
			IncrementTagUsage(tag.ID)
		}
		if len(tags) > 0 {
			database.DB.Model(&topic).Association("Tags").Replace(tags)
		}
	}

	database.DB.Preload("User").Preload("Forum").Preload("Tags").First(&topic, topic.ID)
	utils.Success(w, topic)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("获取话题失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	if topic.UserID != userID {
		log.Printf("更新话题失败: 无权限, topicID=%d, userID=%d", id, userID)
		utils.Error(w, 403, "无权限编辑")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("更新话题请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")
	delete(updates, "is_pinned")
	delete(updates, "is_locked")
	delete(updates, "is_essence")
	delete(updates, "like_count")
	delete(updates, "view_count")
	delete(updates, "reply_count")

	if err := database.DB.Model(&topic).Updates(updates).Error; err != nil {
		log.Printf("更新话题失败: id=%d, userID=%d, 错误: %v", id, userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	database.DB.Preload("User").Preload("Forum").First(&topic, id)
	utils.Success(w, topic)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("删除话题失败: 未认证")
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("获取话题失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("获取用户失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	if topic.UserID != userID && user.Role < 1 {
		log.Printf("删除话题失败: 无权限, topicID=%d, userID=%d, topicUserID=%d", id, userID, topic.UserID)
		utils.Error(w, 403, "无权限删除")
		return
	}

	database.DB.Unscoped().Delete(&topic)
	utils.Success(w, nil)
}
