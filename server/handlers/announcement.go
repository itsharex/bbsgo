package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"net/http"
	"time"
)

func GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	if cached, ok := cache.Get("announcements:list"); ok {
		utils.Success(w, cached)
		return
	}

	var announcements []models.Announcement
	now := time.Now()
	database.DB.Where("(expires_at IS NULL OR expires_at > ?)", now).
		Order("is_pinned DESC, created_at DESC").Find(&announcements)

	cache.Set("announcements:list", announcements, 5*time.Minute)
	utils.Success(w, announcements)
}
