package middleware

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"net/http"
)

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			utils.Error(w, 401, "未认证")
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			utils.Error(w, 401, "用户不存在")
			return
		}

		if user.Role < 2 {
			utils.Error(w, 403, "需要管理员权限")
			return
		}

		next.ServeHTTP(w, r)
	})
}
