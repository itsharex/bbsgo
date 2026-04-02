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

func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.User{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Find(&users)

	utils.Success(w, map[string]interface{}{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req struct {
		Role int `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("更新用户角色请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.Role < 0 || req.Role > 2 {
		log.Printf("更新用户角色失败: 无效的角色, id=%d, role=%d", id, req.Role)
		utils.Error(w, 400, "无效的角色")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("获取用户失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if uint(id) == userID {
		log.Printf("更新用户角色失败: 不能修改自己的角色, id=%d, userID=%d", id, userID)
		utils.Error(w, 400, "不能修改自己的角色")
		return
	}

	database.DB.Model(&user).Update("role", req.Role)
	utils.Success(w, user)
}

func BanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("获取用户失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if uint(id) == userID {
		log.Printf("封禁用户失败: 不能封禁自己, id=%d, userID=%d", id, userID)
		utils.Error(w, 400, "不能封禁自己")
		return
	}

	var req struct {
		Banned bool `json:"banned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("封禁用户请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	database.DB.Model(&user).Update("is_banned", req.Banned)
	utils.Success(w, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	adminID, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("删除用户失败: 未授权, id=%d", id)
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	if uint(id) == adminID {
		log.Printf("删除用户失败: 不能删除自己, id=%d, adminID=%d", id, adminID)
		utils.Error(w, 400, "不能删除自己")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("获取用户失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 开始事务，物理删除用户及其相关数据
	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("开始事务失败: 错误: %v", tx.Error)
		utils.Error(w, 500, "操作失败")
		return
	}

	// 删除用户的帖子 - 使用 Unscoped 物理删除
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Post{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户帖子失败: userId=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除用户帖子失败")
		return
	}

	// 删除用户的收藏
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Favorite{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户收藏失败: userId=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除用户收藏失败")
		return
	}

	// 删除用户的点赞
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Like{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户点赞失败: userId=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除用户点赞失败")
		return
	}

	// 删除用户的通知
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Notification{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户通知失败: userId=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除用户通知失败")
		return
	}

	// 删除用户的消息
	if err := tx.Unscoped().Where("from_user_id = ? OR to_user_id = ?", id, id).Delete(&models.Message{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户消息失败: userId=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除用户消息失败")
		return
	}

	// 删除用户（最后删除用户本身）- 使用 Unscoped 物理删除
	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户失败: userId=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除用户失败")
		return
	}

	tx.Commit()
	log.Printf("删除用户成功: userId=%d, username=%s", id, user.Username)
	utils.Success(w, nil)
}

func CreateForum(w http.ResponseWriter, r *http.Request) {
	var forum models.Forum
	if err := json.NewDecoder(r.Body).Decode(&forum); err != nil {
		log.Printf("创建版块请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if forum.Name == "" {
		log.Printf("创建版块失败: 版块名称为空")
		utils.Error(w, 400, "请填写版块名称")
		return
	}

	if err := database.DB.Create(&forum).Error; err != nil {
		log.Printf("创建版块失败: name=%s, 错误: %v", forum.Name, err)
		utils.Error(w, 500, "创建失败")
		return
	}

	utils.Success(w, forum)
}

func UpdateForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var forum models.Forum
	if err := database.DB.First(&forum, id).Error; err != nil {
		log.Printf("获取版块失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "版块不存在")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("更新版块请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	delete(updates, "id")
	delete(updates, "created_at")

	if err := database.DB.Model(&forum).Updates(updates).Error; err != nil {
		log.Printf("更新版块失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	utils.Success(w, forum)
}

func DeleteForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var forum models.Forum
	if err := database.DB.First(&forum, id).Error; err != nil {
		log.Printf("获取版块失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "版块不存在")
		return
	}

	database.DB.Unscoped().Delete(&forum)
	utils.Success(w, nil)
}

func GetAdminTopics(w http.ResponseWriter, r *http.Request) {
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

	database.DB.Model(&models.Topic{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Preload("User").Preload("Forum").Find(&topics)

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func DeleteAdminTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("获取话题失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	database.DB.Unscoped().Delete(&topic)
	utils.Success(w, nil)
}

func GetAdminPosts(w http.ResponseWriter, r *http.Request) {
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

	database.DB.Model(&models.Post{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Preload("User").Find(&posts)

	utils.Success(w, map[string]interface{}{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func DeleteAdminPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		log.Printf("获取评论失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "评论不存在")
		return
	}

	database.DB.Unscoped().Delete(&post)
	utils.Success(w, nil)
}

func GetAdminReports(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var reports []models.Report
	var total int64

	offset := (page - 1) * pageSize

	database.DB.Model(&models.Report{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Preload("Reporter").Find(&reports)

	utils.Success(w, map[string]interface{}{
		"list":      reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func HandleReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var report models.Report
	if err := database.DB.First(&report, id).Error; err != nil {
		log.Printf("获取举报失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "举报不存在")
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("处理举报请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	now := time.Now()

	database.DB.Model(&report).Updates(map[string]interface{}{
		"status":     req.Status,
		"handled_at": now,
		"handler_id": userID,
	})

	utils.Success(w, report)
}

func CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announcement models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&announcement); err != nil {
		log.Printf("创建公告请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if announcement.Title == "" || announcement.Content == "" {
		log.Printf("创建公告失败: 信息不完整, title=%s", announcement.Title)
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		log.Printf("创建公告失败: title=%s, 错误: %v", announcement.Title, err)
		utils.Error(w, 500, "创建失败")
		return
	}

	utils.Success(w, announcement)
}

func UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		log.Printf("获取公告失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "公告不存在")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("更新公告请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if err := database.DB.Model(&announcement).Updates(updates).Error; err != nil {
		log.Printf("更新公告失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	utils.Success(w, announcement)
}

func DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		log.Printf("获取公告失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "公告不存在")
		return
	}

	database.DB.Unscoped().Delete(&announcement)
	utils.Success(w, nil)
}

func GetAdminTags(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("获取标签列表失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var tags []models.Tag
	database.DB.Order("usage_count DESC, sort_order ASC").Find(&tags)
	utils.Success(w, tags)
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("创建标签失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		Name       string `json:"name"`
		Icon       string `json:"icon"`
		IsOfficial bool   `json:"is_official"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建标签请求解析失败: 错误: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if req.Name == "" {
		log.Printf("创建标签失败: 标签名称为空")
		utils.Error(w, http.StatusBadRequest, "标签名称不能为空")
		return
	}

	tag := models.Tag{
		Name:       req.Name,
		Icon:       req.Icon,
		IsOfficial: req.IsOfficial,
	}
	if err := database.DB.Create(&tag).Error; err != nil {
		log.Printf("创建标签失败: name=%s, 错误: %v", req.Name, err)
		utils.Error(w, 500, "创建标签失败")
		return
	}
	log.Printf("创建标签成功: name=%s", req.Name)
	utils.Success(w, tag)
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("更新标签失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req struct {
		Name       string `json:"name"`
		Icon       string `json:"icon"`
		IsOfficial bool   `json:"is_official"`
		IsBanned   bool   `json:"is_banned"`
		SortOrder  int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("更新标签请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	var tag models.Tag
	if result := database.DB.First(&tag, id); result.Error != nil {
		log.Printf("获取标签失败: id=%d, 错误: %v", id, result.Error)
		utils.Error(w, http.StatusNotFound, "话题标签不存在")
		return
	}

	tag.Name = req.Name
	tag.Icon = req.Icon
	tag.IsOfficial = req.IsOfficial
	tag.IsBanned = req.IsBanned
	tag.SortOrder = req.SortOrder

	if err := database.DB.Save(&tag).Error; err != nil {
		log.Printf("更新标签失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 500, "更新标签失败")
		return
	}
	log.Printf("更新标签成功: id=%d, name=%s", id, req.Name)
	utils.Success(w, tag)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("删除标签失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := database.DB.Unscoped().Delete(&models.Tag{}, id).Error; err != nil {
		log.Printf("删除标签失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除标签失败")
		return
	}
	log.Printf("删除标签成功: id=%d", id)
	utils.Success(w, nil)
}

func MergeTags(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("合并标签失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		SourceID uint `json:"source_id"`
		TargetID uint `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("合并标签请求解析失败: 错误: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if req.SourceID == 0 || req.TargetID == 0 || req.SourceID == req.TargetID {
		log.Printf("合并标签失败: 无效的标签ID, sourceID=%d, targetID=%d", req.SourceID, req.TargetID)
		utils.Error(w, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var sourceTag, targetTag models.Tag
	if err := database.DB.First(&sourceTag, req.SourceID).Error; err != nil {
		log.Printf("获取源标签失败: sourceID=%d, 错误: %v", req.SourceID, err)
		utils.Error(w, http.StatusNotFound, "源标签不存在")
		return
	}
	if err := database.DB.First(&targetTag, req.TargetID).Error; err != nil {
		log.Printf("获取目标标签失败: targetID=%d, 错误: %v", req.TargetID, err)
		utils.Error(w, http.StatusNotFound, "目标标签不存在")
		return
	}

	var topics []models.Topic
	database.DB.Model(&sourceTag).Association("Topics").Find(&topics)

	if len(topics) > 0 {
		database.DB.Model(&targetTag).Association("Topics").Append(topics)
		database.DB.Model(&sourceTag).Association("Topics").Clear()
	}

	targetTag.UsageCount += sourceTag.UsageCount
	database.DB.Save(&targetTag)

	database.DB.Unscoped().Delete(&sourceTag)
	log.Printf("合并标签成功: sourceID=%d, targetID=%d, mergedCount=%d", req.SourceID, req.TargetID, len(topics))
	utils.Success(w, map[string]interface{}{
		"merged_count": len(topics),
		"target_tag":   targetTag,
	})
}

func ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("修改密码失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("修改密码请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		log.Printf("修改密码失败: 信息不完整")
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	if len(req.NewPassword) < 6 {
		log.Printf("修改密码失败: 新密码长度不足, length=%d", len(req.NewPassword))
		utils.Error(w, 400, "新密码长度至少为6位")
		return
	}

	var admin models.User
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		log.Printf("获取管理员失败: id=%d, 错误: %v", adminID, err)
		utils.Error(w, 404, "管理员不存在")
		return
	}

	if !utils.CheckPassword(req.OldPassword, admin.PasswordHash) {
		log.Printf("修改密码失败: 原密码错误, adminID=%d", adminID)
		utils.Error(w, 400, "原密码错误")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("密码加密失败: 错误: %v", err)
		utils.Error(w, 500, "密码加密失败")
		return
	}

	if err := database.DB.Model(&admin).Update("password_hash", hashedPassword).Error; err != nil {
		log.Printf("更新密码失败: adminID=%d, 错误: %v", adminID, err)
		utils.Error(w, 500, "更新密码失败")
		return
	}
	log.Printf("修改密码成功: adminID=%d", adminID)
	utils.Success(w, nil)
}
