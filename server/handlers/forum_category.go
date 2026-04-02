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

func GetForumCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.ForumCategory
	database.DB.Where("is_active = ?", true).Order("sort_order").Find(&categories)
	utils.Success(w, categories)
}

func GetAllForumCategories(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("获取所有分类失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var categories []models.ForumCategory
	database.DB.Order("sort_order").Find(&categories)
	utils.Success(w, categories)
}

func CreateForumCategory(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("创建分类失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req models.ForumCategory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("创建分类请求解析失败: 错误: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		log.Printf("创建分类失败: name=%s, 错误: %v", req.Name, err)
		utils.Error(w, 500, "创建分类失败")
		return
	}
	log.Printf("创建分类成功: name=%s", req.Name)
	utils.Success(w, req)
}

func UpdateForumCategory(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("更新分类失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req models.ForumCategory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("更新分类请求解析失败: id=%d, 错误: %v", id, err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	var category models.ForumCategory
	if result := database.DB.First(&category, id); result.Error != nil {
		log.Printf("获取分类失败: id=%d, 错误: %v", id, result.Error)
		utils.Error(w, http.StatusNotFound, "分类不存在")
		return
	}

	category.Name = req.Name
	category.Icon = req.Icon
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	category.IsActive = req.IsActive

	if err := database.DB.Save(&category).Error; err != nil {
		log.Printf("更新分类失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 500, "更新分类失败")
		return
	}
	log.Printf("更新分类成功: id=%d, name=%s", id, req.Name)
	utils.Success(w, category)
}

func DeleteForumCategory(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("删除分类失败: 未授权")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := database.DB.Unscoped().Delete(&models.ForumCategory{}, id).Error; err != nil {
		log.Printf("删除分类失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 500, "删除分类失败")
		return
	}
	log.Printf("删除分类成功: id=%d", id)
	utils.Success(w, nil)
}
