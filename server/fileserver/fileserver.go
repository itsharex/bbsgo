package fileserver

import (
	"embed"
	"log"
	"net/http"
	"path"
)

// AdminFS 管理后台静态文件的 embed 实例
//
//go:embed admin
var AdminFS embed.FS

// SiteFS 主站静态文件的 embed 实例
//
//go:embed site
var SiteFS embed.FS

// ServeAdmin 处理管理后台静态文件
// mux 会把 /console/xxx 路径strip成 /xxx 传给我们
func ServeAdmin(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	// 去掉 /console 前缀
	if len(reqPath) >= 8 && reqPath[:8] == "/console" {
		reqPath = reqPath[8:] // /console/xxx -> /xxx
		if reqPath == "" {
			reqPath = "/"
		}
	}
	log.Printf("[admin] serving: %s", reqPath)
	serveFile(w, r, reqPath, AdminFS, "admin")
}

// ServeSite 处理主站静态文件
func ServeSite(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	log.Printf("[site] serving: %s", reqPath)
	serveFile(w, r, reqPath, SiteFS, "site")
}

func serveFile(w http.ResponseWriter, r *http.Request, reqPath string, fsys embed.FS, name string) {
	// 清理路径
	reqPath = path.Clean(reqPath)
	if reqPath == "" || reqPath == "." {
		reqPath = "/"
	}

	// 拼接实际的文件路径: admin + /index.html
	filePath := name + reqPath

	// 读取文件
	content, err := fsys.ReadFile(filePath)
	if err == nil {
		log.Printf("[%s] found: %s, size: %d", name, filePath, len(content))
		w.Header().Set("Content-Type", contentType(reqPath))
		w.Write(content)
		return
	}

	// 文件不存在，尝试 index.html
	indexPath := name + "/index.html"
	content, err = fsys.ReadFile(indexPath)
	if err != nil {
		log.Printf("[%s] not found: %s, index.html not found: %v", name, filePath, err)
		http.Error(w, "Not Found", 404)
		return
	}

	log.Printf("[%s] serving index.html for: %s, size: %d", name, reqPath, len(content))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// contentType 根据文件扩展名返回 Content-Type
func contentType(filePath string) string {
	ext := path.Ext(filePath)
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
