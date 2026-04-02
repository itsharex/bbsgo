package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"net/http"
	"time"
)

func GetForums(w http.ResponseWriter, r *http.Request) {
	if cached, ok := cache.Get("forums:list"); ok {
		utils.Success(w, cached)
		return
	}

	var forums []models.Forum
	database.DB.Order("sort_order ASC, id ASC").Find(&forums)

	cache.Set("forums:list", forums, 5*time.Minute)
	utils.Success(w, forums)
}
