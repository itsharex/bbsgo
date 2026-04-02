package handlers

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"
)

func Search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		log.Printf("搜索失败: 关键词为空")
		utils.Error(w, 400, "请输入搜索关键词")
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
	searchPattern := "%" + keyword + "%"

	database.DB.Model(&models.Topic{}).Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).Count(&total)
	database.DB.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Preload("User").Preload("Forum").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&topics)

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"keyword":   keyword,
	})
}
