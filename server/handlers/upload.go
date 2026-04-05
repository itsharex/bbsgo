package handlers

import (
	"bbsgo/errors"
	"bbsgo/storage"
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// CheckFileExists 检查文件是否已存在（用于秒传）
// 通过文件内容hash生成相同key，如果文件存在则直接返回URL
func CheckFileExists(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	contentHash := r.URL.Query().Get("content_hash") // 文件内容MD5
	if filename == "" {
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}
	if contentHash == "" {
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 获取存储服务实例
	storageSvc, err := storage.GetStorage()
	if err != nil {
		log.Printf("[upload/exists] failed to get storage service, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 根据文件扩展名确定目录
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true}
	videoExts := map[string]bool{".mp4": true, ".webm": true, ".ogg": true, ".mov": true, ".mkv": true, ".avi": true}

	dir := ""
	if imageExts[ext] {
		dir = "images"
	} else if videoExts[ext] {
		dir = "videos"
	}

	// 生成文件key（使用content_hash确保相同内容生成相同key）
	key := storage.GenerateFileKeyWithHash(dir, filename, contentHash)

	// 检查文件是否存在
	if storageSvc.Exists(key) {
		url := storageSvc.GetURL(key)
		log.Printf("[upload/exists] file exists, key: %s", key)
		errors.Success(w, map[string]interface{}{
			"exists": true,
			"url":    url,
			"key":    key,
		})
		return
	}

	log.Printf("[upload/exists] file not exists, key: %s", key)
	errors.Success(w, map[string]interface{}{
		"exists": false,
		"key":    key,
	})
}

// validateFileType 验证文件类型（通过文件内容）
func validateFileType(fileData []byte, ext string) (bool, string) {
	// 检测文件的真实类型
	contentType := http.DetectContentType(fileData)
	
	// 支持的图片类型
	imageTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"image/bmp":  true,
	}
	
	// 支持的视频类型
	videoTypes := map[string]bool{
		"video/mp4":        true,
		"video/webm":       true,
		"video/ogg":        true,
		"video/quicktime":  true,
		"video/x-msvideo":  true,
		"video/x-matroska": true,
	}
	
	// 检查是否为图片
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
	}
	
	// 检查是否为视频
	videoExts := map[string]bool{
		".mp4":  true,
		".webm": true,
		".ogg":  true,
		".mov":  true,
		".mkv":  true,
		".avi":  true,
	}
	
	// 验证扩展名和内容类型是否匹配
	if imageExts[ext] {
		if imageTypes[contentType] {
			return true, contentType
		}
		// 某些图片可能被识别为 application/octet-stream
		if contentType == "application/octet-stream" {
			return true, "image/jpeg"
		}
		return false, contentType
	}
	
	if videoExts[ext] {
		if videoTypes[contentType] {
			return true, contentType
		}
		// 某些视频可能被识别为 application/octet-stream
		if contentType == "application/octet-stream" {
			return true, "video/mp4"
		}
		return false, contentType
	}
	
	return false, contentType
}

// checkMaliciousContent 检查文件是否包含恶意内容
func checkMaliciousContent(fileData []byte, ext string) bool {
	// SVG 文件可能包含恶意脚本
	if ext == ".svg" {
		// 检查是否包含脚本标签
		content := strings.ToLower(string(fileData))
		if strings.Contains(content, "<script") ||
			strings.Contains(content, "javascript:") ||
			strings.Contains(content, "onerror=") ||
			strings.Contains(content, "onload=") ||
			strings.Contains(content, "onclick=") {
			return true
		}
	}
	
	// 检查图片文件是否包含 PHP/ASP 代码（防止图片马）
	if ext != ".svg" {
		content := strings.ToLower(string(fileData))
		if strings.Contains(content, "<?php") ||
			strings.Contains(content, "<%") ||
			strings.Contains(content, "<script") {
			return true
		}
	}
	
	return false
}

// UploadFile 文件上传处理器
// 支持图片和视频上传到配置的存储服务（本地/七牛云/阿里云/腾讯云）
// 最大文件大小：图片50MB，视频500MB
func UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[upload] upload handler started")

	// 解析 multipart 表单，最大500MB
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		log.Printf("[upload] failed to parse multipart form, error: %v", err)
		errors.Error(w, errors.CodeFileSizeExceeded, "")
		return
	}

	// 获取上传文件
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[upload] failed to get form file, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}
	defer file.Close()

	log.Printf("[upload] received file: %s, size: %d bytes", header.Filename, header.Size)

	// 获取文件扩展名并验证
	ext := strings.ToLower(filepath.Ext(header.Filename))

	// 支持的文件格式（移除 SVG）
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
	}

	videoExts := map[string]bool{
		".mp4":  true,
		".webm": true,
		".ogg":  true,
		".mov":  true,
		".mkv":  true,
		".avi":  true,
	}

	// 验证文件类型
	isImage := imageExts[ext]
	isVideo := videoExts[ext]

	if !isImage && !isVideo {
		log.Printf("[upload] unsupported file type: %s", ext)
		errors.Error(w, errors.CodeFileTypeUnsupported, "")
		return
	}

	// 验证文件大小（图片50MB，视频500MB）
	if isImage && header.Size > 50*1024*1024 {
		log.Printf("[upload] image too large: %d bytes", header.Size)
		errors.Error(w, errors.CodeImageSizeExceeded, "")
		return
	}

	if isVideo && header.Size > 500*1024*1024 {
		log.Printf("[upload] video too large: %d bytes", header.Size)
		errors.Error(w, errors.CodeFileSizeExceeded, "")
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[upload] failed to read file data, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}
	log.Printf("[upload] file read success, size: %d bytes", len(fileData))

	// 验证文件内容类型
	valid, detectedType := validateFileType(fileData, ext)
	if !valid {
		log.Printf("[upload] file content type mismatch, ext: %s, detected: %s", ext, detectedType)
		errors.Error(w, errors.CodeFileTypeUnsupported, "File content does not match extension")
		return
	}

	// 检查恶意内容
	if checkMaliciousContent(fileData, ext) {
		log.Printf("[upload] malicious content detected in file: %s", header.Filename)
		errors.Error(w, errors.CodeFileTypeUnsupported, "Malicious content detected")
		return
	}

	// 获取存储服务实例
	storageSvc, err := storage.GetStorage()
	if err != nil {
		log.Printf("[upload] failed to get storage service, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 获取上传目录参数
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		if isImage {
			dir = "images"
		} else if isVideo {
			dir = "videos"
		}
	}

	// 获取content_hash参数（可选，用于秒传）
	contentHash := r.URL.Query().Get("content_hash")

	// 生成存储文件key
	// 如果提供了content_hash，使用GenerateFileKeyWithHash确保与秒传检查一致
	var key string
	if contentHash != "" {
		key = storage.GenerateFileKeyWithHash(dir, header.Filename, contentHash)
		log.Printf("[upload] using content hash for key")
	} else {
		key = storage.GenerateFileKey(dir, header.Filename)
		log.Printf("[upload] using filename hash for key")
	}

	contentType := detectedType
	if contentType == "" {
		contentType = header.Header.Get("Content-Type")
		if contentType == "" {
			if isImage {
				contentType = "image/jpeg"
			} else if isVideo {
				contentType = "video/mp4"
			} else {
				contentType = "application/octet-stream"
			}
		}
	}

	// 上传到存储服务
	url, err := storageSvc.Upload(key, fileData, contentType)
	if err != nil {
		log.Printf("[upload] failed to upload to storage, error: %v", err)
		errors.Error(w, errors.CodeUploadFailed, "")
		return
	}

	log.Printf("[upload] upload success, key: %s", key)
	log.Printf("[upload] upload handler finished")

	errors.Success(w, map[string]string{
		"url": url,
	})
}

// SanitizeSVG 清理 SVG 文件中的潜在危险内容
func SanitizeSVG(fileData []byte) []byte {
	content := string(fileData)
	
	// 移除 script 标签
	content = strings.ReplaceAll(content, "<script", "&lt;script")
	content = strings.ReplaceAll(content, "</script>", "&lt;/script&gt;")
	
	// 移除事件处理器
	eventHandlers := []string{"onerror", "onload", "onclick", "onmouseover", "onfocus", "onblur"}
	for _, handler := range eventHandlers {
		content = strings.ReplaceAll(content, handler+"=", "data-removed=")
	}
	
	// 移除 javascript: 伪协议
	content = strings.ReplaceAll(content, "javascript:", "")
	
	return []byte(content)
}

// UploadSVG SVG 文件上传处理器（特殊处理）
func UploadSVG(w http.ResponseWriter, r *http.Request) {
	log.Printf("[upload/svg] svg upload handler started")

	// 解析 multipart 表单，最大10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("[upload/svg] failed to parse multipart form, error: %v", err)
		errors.Error(w, errors.CodeFileSizeExceeded, "")
		return
	}

	// 获取上传文件
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[upload/svg] failed to get form file, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}
	defer file.Close()

	// 验证扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".svg" {
		log.Printf("[upload/svg] invalid file extension: %s", ext)
		errors.Error(w, errors.CodeFileTypeUnsupported, "")
		return
	}

	// 验证文件大小（最大10MB）
	if header.Size > 10*1024*1024 {
		log.Printf("[upload/svg] file too large: %d bytes", header.Size)
		errors.Error(w, errors.CodeImageSizeExceeded, "")
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[upload/svg] failed to read file data, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 检测内容类型
	contentType := http.DetectContentType(fileData)
	if contentType != "image/svg+xml" && contentType != "text/xml" && contentType != "application/xml" {
		// SVG 可能被识别为 text/plain 或 application/octet-stream
		if !bytes.Contains(fileData, []byte("<svg")) {
			log.Printf("[upload/svg] invalid svg content type: %s", contentType)
			errors.Error(w, errors.CodeFileTypeUnsupported, "Invalid SVG file")
			return
		}
	}

	// 检查恶意内容
	if checkMaliciousContent(fileData, ".svg") {
		log.Printf("[upload/svg] malicious content detected")
		errors.Error(w, errors.CodeFileTypeUnsupported, "Malicious content detected in SVG")
		return
	}

	// 清理 SVG 内容
	sanitizedData := SanitizeSVG(fileData)

	// 获取存储服务实例
	storageSvc, err := storage.GetStorage()
	if err != nil {
		log.Printf("[upload/svg] failed to get storage service, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 生成存储文件key
	key := storage.GenerateFileKey("images", header.Filename)

	// 上传到存储服务
	url, err := storageSvc.Upload(key, sanitizedData, "image/svg+xml")
	if err != nil {
		log.Printf("[upload/svg] failed to upload to storage, error: %v", err)
		errors.Error(w, errors.CodeUploadFailed, "")
		return
	}

	log.Printf("[upload/svg] upload success, key: %s", key)

	errors.Success(w, map[string]string{
		"url": url,
	})
}
