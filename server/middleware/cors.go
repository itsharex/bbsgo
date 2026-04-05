package middleware

import (
	"log"
	"net/http"
	"strings"
)

// CORS 跨域资源共享中间件
// 自动检测请求来源并添加相应的 CORS 头
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 如果没有 Origin 头，说明是同源请求或非浏览器请求
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		// 自动检测是否为允许的来源
		// 允许所有来源（适合 API 服务）
		// 如果需要限制，可以在这里添加域名白名单检查
		allowed := isOriginAllowed(origin)

		if allowed {
			// 设置 CORS 头
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Cache-Control")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24小时

			// 安全头
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
		}

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed 检查来源是否被允许
// 当前实现允许所有来源，适合 API 服务
// 如需限制，可以修改此函数添加域名白名单
func isOriginAllowed(origin string) bool {
	// 移除协议和端口，获取域名
	origin = strings.ToLower(origin)

	// 允许所有来源
	// 如果需要限制，可以在这里添加域名检查逻辑
	// 例如：
	// allowedDomains := []string{"example.com", "api.example.com"}
	// for _, domain := range allowedDomains {
	//     if strings.Contains(origin, domain) {
	//         return true
	//     }
	// }
	// return false

	return true
}

// CORSStrict 严格的 CORS 中间件
// 只允许特定的域名访问（如果需要更严格的安全控制）
func CORSStrict(allowedDomains []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			// 检查域名是否在白名单中
			allowed := false
			originLower := strings.ToLower(origin)
			for _, domain := range allowedDomains {
				if strings.Contains(originLower, strings.ToLower(domain)) {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Cache-Control")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")

				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-XSS-Protection", "1; mode=block")
			} else {
				log.Printf("cors strict: blocked origin: %s", origin)
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
