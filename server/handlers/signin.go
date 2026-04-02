package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"time"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("获取用户失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	today := time.Now().Format("2006-01-02")
	if user.LastSignAt != nil && user.LastSignAt.Format("2006-01-02") == today {
		log.Printf("签到失败: 今日已签到, userID=%d", userID)
		utils.Error(w, 400, "今日已签到")
		return
	}

	credits := 10

	var lastSign time.Time
	if user.LastSignAt != nil {
		lastSign = *user.LastSignAt
	}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if lastSign.Format("2006-01-02") == yesterday {
		credits = 15
	}

	now := time.Now()
	user.LastSignAt = &now
	user.Credits += credits

	database.DB.Save(&user)

	utils.Success(w, map[string]interface{}{
		"credits":      credits,
		"total_credits": user.Credits,
		"message":      "签到成功",
	})
}

func GetSignInStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("获取用户失败: userID=%d, 错误: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	today := time.Now().Format("2006-01-02")
	signedToday := user.LastSignAt != nil && user.LastSignAt.Format("2006-01-02") == today

	utils.Success(w, map[string]interface{}{
		"signed_today": signedToday,
		"last_sign_at": user.LastSignAt,
		"credits":      user.Credits,
	})
}
