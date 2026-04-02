package middleware

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, 401, "未提供认证令牌")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Error(w, 401, "无效的认证令牌")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			utils.Error(w, 401, "未授权")
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

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	claims, ok := ctx.Value(UserContextKey).(*utils.Claims)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

func GetAdminIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return 0, false
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return 0, false
	}

	if user.Role < 2 {
		return 0, false
	}

	return userID, true
}
