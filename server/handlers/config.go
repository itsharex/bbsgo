package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"net/http"
)

func GetSiteConfig(w http.ResponseWriter, r *http.Request) {
	var configs []models.SiteConfig
	database.DB.Find(&configs)

	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	utils.Success(w, configMap)
}

func UpdateSiteConfig(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	for key, value := range req {
		var config models.SiteConfig
		result := database.DB.Where("key = ?", key).First(&config)
		if result.Error == nil {
			config.Value = value
			database.DB.Save(&config)
		} else {
			newConfig := models.SiteConfig{
				Key:   key,
				Value: value,
			}
			database.DB.Create(&newConfig)
		}
	}

	utils.Success(w, nil)
}
