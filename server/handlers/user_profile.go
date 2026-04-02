package handlers

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("获取用户失败: 无效的用户ID, id=%s", vars["id"])
		utils.Error(w, 400, "无效的用户ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("获取用户失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	utils.Success(w, user)
}

func GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("获取用户统计失败: 无效的用户ID, id=%s", vars["id"])
		utils.Error(w, 400, "无效的用户ID")
		return
	}

	var topicCount int64
	var postCount int64
	var rank int64

	database.DB.Model(&models.Topic{}).Where("user_id = ?", id).Count(&topicCount)
	database.DB.Model(&models.Post{}).Where("user_id = ?", id).Count(&postCount)
	database.DB.Table("users").Where("id <= ?", id).Count(&rank)

	utils.Success(w, map[string]interface{}{
		"topic_count": topicCount,
		"post_count":  postCount,
		"rank":        rank,
	})
}

func GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("获取粉丝列表失败: 无效的用户ID, id=%s", vars["id"])
		utils.Error(w, 400, "无效的用户ID")
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

	var follows []models.Follow
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", id).Count(&total)
	database.DB.Where("follow_user_id = ?", id).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&follows)

	followers := make([]models.User, 0)
	for _, follow := range follows {
		followers = append(followers, follow.User)
	}

	utils.Success(w, map[string]interface{}{
		"list":      followers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetUserTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("获取用户帖子失败: 无效的用户ID, id=%s", vars["id"])
		utils.Error(w, 400, "无效的用户ID")
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

	database.DB.Model(&models.Topic{}).Where("user_id = ?", id).Count(&total)
	database.DB.Where("user_id = ?", id).
		Preload("User").
		Preload("Forum").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&topics)

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
