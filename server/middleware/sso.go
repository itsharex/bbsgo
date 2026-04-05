package middleware

import (
	"bbsgo/services"
	"bbsgo/utils"
	"context"
	"log"
	"net/http"
)

// SSOCookieNameContextKey SSO cookie 名上下文键
type SSOCookieNameContextKey string

// SSOTokenContextKey SSO token 上下文键（已验证的主站 token）
const SSOTokenContextKey = SSOCookieNameContextKey("sso_token")

// SSOMiddleware SSO 自动登录中间件
// 在每个 API 请求时检查共享域 cookie 中的主站 token，自动完成 SSO 登录
// 登录成功后生成 BBSGO token 并写入 response cookie
func SSOMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// OPTIONS 请求直接放行
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// 如果已经通过 Authorization header 认证，跳过 SSO
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			next.ServeHTTP(w, r)
			return
		}

		// 获取 SSO 配置
		cookieName := utils.GetConfigString("sso_cookie_name", "token")
		verifyURL := utils.GetConfigString("sso_verify_url", "")
		if verifyURL == "" {
			// 未配置 SSO，跳过
			next.ServeHTTP(w, r)
			return
		}

		// 从 cookie 获取主站 token
		token := ""
		if cookie, err := r.Cookie(cookieName); err == nil {
			token = cookie.Value
		}

		if token == "" {
			// 没有 SSO cookie，跳过
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("[sso middleware] found token in cookie %s", cookieName)

		// 验证主站 token 并获取用户信息
		userInfo, err := services.VerifySSO(verifyURL, token)
		if err != nil {
			log.Printf("[sso middleware] verify failed: %v", err)
			// 验证失败，但不影响正常流程，继续让后续中间件/handler 处理
			next.ServeHTTP(w, r)
			return
		}

		// 查找或创建本地用户
		user, err := services.GetOrCreateSSOUser(userInfo)
		if err != nil {
			log.Printf("[sso middleware] get or create user failed: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		// 生成 BBSGO token
		bbsgoToken, err := utils.GenerateToken(user.ID, user.Username, user.TokenVersion)
		if err != nil {
			log.Printf("[sso middleware] generate token failed: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		// 将 BBSGO token 写入 response cookie（HttpOnly: false 让前端 JS 可读）
		http.SetCookie(w, &http.Cookie{
			Name:     "bbsgo_token",
			Value:    bbsgoToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})

		// 在当前请求的 header 中设置 Authorization，让当前请求也能使用
		r.Header.Set("Authorization", "Bearer "+bbsgoToken)

		// 将用户信息存入 context
		claims := &utils.Claims{
			UserID:   user.ID,
			Username: user.Username,
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		ctx = context.WithValue(ctx, SSOTokenContextKey, token)

		log.Printf("[sso middleware] auto login success: user_id=%d, username=%s", user.ID, user.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
