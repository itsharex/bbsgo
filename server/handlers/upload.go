package handlers

import (
	"bbsgo/services"
	"bbsgo/utils"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Upload] ========== UploadFileHandler Start ==========")

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Printf("[Upload] ERROR: Parse multipart form failed: %v", err)
		utils.Error(w, 400, "文件大小超过限制(最大50MB)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[Upload] ERROR: Get form file failed: %v", err)
		utils.Error(w, 400, "获取文件失败")
		return
	}
	defer file.Close()

	log.Printf("[Upload] Received file: %s, Size: %d bytes", header.Filename, header.Size)

	if header.Size > 50*1024*1024 {
		log.Printf("[Upload] ERROR: File too large: %d bytes", header.Size)
		utils.Error(w, 400, "文件大小超过限制(最大50MB)")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		log.Printf("[Upload] ERROR: Unsupported file type: %s", ext)
		utils.Error(w, 400, "不支持的文件类型，仅支持 jpg, jpeg, png, gif, webp")
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[Upload] ERROR: Read file failed: %v", err)
		utils.Error(w, 500, "读取文件失败")
		return
	}
	log.Printf("[Upload] File read success: %d bytes", len(fileData))

	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "uploads"
	}

	url, err := services.UploadToQiniu(fileData, header.Filename, dir)
	if err != nil {
		log.Printf("[Upload] ERROR: Upload to Qiniu failed: %v", err)
		utils.Error(w, 500, "上传到七牛云失败")
		return
	}

	log.Printf("[Upload] Upload success, URL: %s", url)
	log.Printf("[Upload] ========== UploadFileHandler End ==========")

	utils.Success(w, map[string]string{
		"url": url,
	})
}
