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

func GetFollows(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var follows []models.Follow
	database.DB.Where("user_id = ?", userID).Preload("FollowUser").Find(&follows)

	utils.Success(w, follows)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var followers []models.Follow
	database.DB.Where("follow_user_id = ?", userID).Preload("User").Find(&followers)

	utils.Success(w, followers)
}

func CreateFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var req struct {
		FollowUserID uint `json:"follow_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建关注请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.FollowUserID == userID {
		log.Printf("关注用户失败: 不能关注自己, userID=%d", userID)
		utils.Error(w, 400, "不能关注自己")
		return
	}

	var followUser models.User
	if err := database.DB.First(&followUser, req.FollowUserID).Error; err != nil {
		log.Printf("获取被关注用户失败: userID=%d, 错误: %v", req.FollowUserID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	var existing models.Follow
	if err := database.DB.Where("user_id = ? AND follow_user_id = ?", userID, req.FollowUserID).First(&existing).Error; err == nil {
		log.Printf("关注用户失败: 已关注, userID=%d, followUserID=%d", userID, req.FollowUserID)
		utils.Error(w, 400, "已关注该用户")
		return
	}

	follow := models.Follow{
		UserID:       userID,
		FollowUserID: req.FollowUserID,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		log.Printf("创建关注失败: userID=%d, followUserID=%d, 错误: %v", userID, req.FollowUserID, err)
		utils.Error(w, 500, "关注失败")
		return
	}

	utils.Success(w, follow)
}

func DeleteFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var req struct {
		FollowUserID uint `json:"follow_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("取消关注请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	database.DB.Unscoped().Where("user_id = ? AND follow_user_id = ?", userID, req.FollowUserID).Delete(&models.Follow{})

	utils.Success(w, nil)
}

func CheckFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	targetUserID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	var follow models.Follow
	isFollowing := database.DB.Where("user_id = ? AND follow_user_id = ?", userID, targetUserID).First(&follow).Error == nil

	utils.Success(w, map[string]bool{"is_following": isFollowing})
}
