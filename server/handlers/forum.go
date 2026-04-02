package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"time"
)

// GetForums 获取版块列表处理器
// 支持缓存，缓存时间5分钟
// 注意：过滤掉"全部"这个虚拟板块，它仅用于首页展示
func GetForums(w http.ResponseWriter, r *http.Request) {
	// 尝试从缓存获取
	if cached, ok := cache.Get("forums:list"); ok {
		utils.Success(w, cached)
		return
	}

	// 查询数据库，排除"全部"板块
	var forums []models.Forum
	if err := database.DB.Where("name != ?", "全部").Order("sort_order ASC, id ASC").Find(&forums).Error; err != nil {
		log.Printf("get forums: failed to query forums, error: %v", err)
		utils.Error(w, 500, "获取版块列表失败")
		return
	}

	// 设置缓存
	cache.Set("forums:list", forums, 5*time.Minute)

	utils.Success(w, forums)
}
