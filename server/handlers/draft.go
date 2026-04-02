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

func GetDrafts(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var drafts []models.Draft
	database.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&drafts)

	utils.Success(w, drafts)
}

func GetDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&draft).Error; err != nil {
		log.Printf("获取草稿失败: id=%d, userID=%d, 错误: %v", id, userID, err)
		utils.Error(w, 404, "草稿不存在")
		return
	}

	utils.Success(w, draft)
}

func CreateDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var draft models.Draft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		log.Printf("创建草稿请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	draft.UserID = userID

	if err := database.DB.Create(&draft).Error; err != nil {
		log.Printf("创建草稿失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 500, "保存失败")
		return
	}

	utils.Success(w, draft)
}

func UpdateDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&draft).Error; err != nil {
		log.Printf("获取草稿失败: id=%d, userID=%d, 错误: %v", id, userID, err)
		utils.Error(w, 404, "草稿不存在")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("更新草稿请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")

	if err := database.DB.Model(&draft).Updates(updates).Error; err != nil {
		log.Printf("更新草稿失败: id=%d, userID=%d, 错误: %v", id, userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	utils.Success(w, draft)
}

func DeleteDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&draft).Error; err != nil {
		log.Printf("获取草稿失败: id=%d, userID=%d, 错误: %v", id, userID, err)
		utils.Error(w, 404, "草稿不存在")
		return
	}

	database.DB.Unscoped().Delete(&draft)
	utils.Success(w, nil)
}
