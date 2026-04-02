package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func GetTags(w http.ResponseWriter, r *http.Request) {
	if cached, ok := cache.Get("tags:hot"); ok {
		utils.Success(w, cached)
		return
	}

	var tags []models.Tag
	database.DB.Where("is_banned = ?", false).
		Order("is_official DESC, usage_count DESC, sort_order ASC").
		Limit(20).
		Find(&tags)

	cache.Set("tags:hot", tags, 10*time.Minute)
	utils.Success(w, tags)
}

func GetTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		log.Printf("获取标签失败: id=%d, 错误: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	utils.Success(w, tag)
}

func SearchTags(w http.ResponseWriter, r *http.Request) {
	keyword := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(keyword) < 1 {
		utils.Success(w, []models.Tag{})
		return
	}

	var tags []models.Tag
	database.DB.Where("is_banned = ? AND name LIKE ?", false, "%"+keyword+"%").
		Order("usage_count DESC").
		Limit(10).
		Find(&tags)

	utils.Success(w, tags)
}

func GetOrCreateTagByName(name string) (*models.Tag, error) {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 20 {
		return nil, nil
	}

	var tag models.Tag
	if err := database.DB.Where("name = ?", name).First(&tag).Error; err != nil {
		tag = models.Tag{
			Name:       name,
			UsageCount: 0,
			IsOfficial: false,
			IsBanned:   false,
		}
		if err := database.DB.Create(&tag).Error; err != nil {
			return nil, err
		}
	}
	return &tag, nil
}

func IncrementTagUsage(tagID uint) {
	database.DB.Model(&models.Tag{}).Where("id = ?", tagID).
		UpdateColumn("usage_count", database.DB.Raw("usage_count + 1"))
}
