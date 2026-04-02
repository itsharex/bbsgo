package middleware

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// OPTIONS 请求直接放行，让 CORS 预检通过
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Auth中间件: 未提供认证令牌, path=%s, method=%s", r.URL.Path, r.Method)
			utils.Error(w, 401, "未提供认证令牌")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("Auth中间件: 解析Token失败, path=%s, method=%s, error=%v", r.URL.Path, r.Method, err)
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
			log.Printf("AdminAuth中间件: 未从上下文获取到用户ID, path=%s, method=%s", r.URL.Path, r.Method)
			utils.Error(w, 401, "未授权")
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			log.Printf("AdminAuth中间件: 用户不存在, userID=%d, path=%s, error=%v", userID, r.URL.Path, err)
			utils.Error(w, 401, "用户不存在")
			return
		}

		if user.Role < 2 {
			log.Printf("AdminAuth中间件: 用户权限不足, userID=%d, role=%d, path=%s", userID, user.Role, r.URL.Path)
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
