package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// 固定上传根目录
const localRootPath = "./uploads"

// LocalStorage 本地存储服务
// 将文件存储在本地文件系统
type LocalStorage struct{}

// newLocalStorage 创建本地存储服务实例
func newLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// NewLocalStorageWithCheck 创建本地存储并检查配置
func NewLocalStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newLocalStorage()
	// 确保上传目录存在
	if err := os.MkdirAll(localRootPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create root directory: %v", err)
	}
	return storage, nil
}

// Name 返回存储类型名称
func (s *LocalStorage) Name() string {
	return "local"
}

// isPathSafe 验证路径是否安全，不允许超出上传根目录
// 返回: 完整安全路径和错误
func isPathSafe(key string) (string, error) {
	// 清理 key 中的 ../ 和 ../
	key = strings.ReplaceAll(key, "..", "")

	// 解析完整路径
	fullPath := filepath.Join(localRootPath, key)

	// 解析并获取绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %v", err)
	}

	// 获取根目录的绝对路径
	absRoot, err := filepath.Abs(localRootPath)
	if err != nil {
		return "", fmt.Errorf("invalid root path: %v", err)
	}

	// 验证路径在根目录内
	if !strings.HasPrefix(absPath, absRoot+string(filepath.Separator)) && absPath != absRoot {
		return "", fmt.Errorf("path outside root directory: %s", key)
	}

	return absPath, nil
}

// Upload 上传文件到本地
// key: 文件存储路径（相对于基础路径）
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *LocalStorage) Upload(key string, data []byte, contentType string) (string, error) {
	// 安全检查路径
	safePath, err := isPathSafe(key)
	if err != nil {
		return "", fmt.Errorf("unsafe path: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(safePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(safePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return s.GetURL(key), nil
}

// Delete 删除本地文件
// key: 文件存储路径
// 返回: 错误
func (s *LocalStorage) Delete(key string) error {
	// 安全检查路径
	safePath, err := isPathSafe(key)
	if err != nil {
		return fmt.Errorf("unsafe path: %v", err)
	}

	return os.Remove(safePath)
}

// Exists 检查本地文件是否存在
// key: 文件存储路径
// 返回: 是否存在
func (s *LocalStorage) Exists(key string) bool {
	safePath, err := isPathSafe(key)
	if err != nil {
		return false
	}
	_, err = os.Stat(safePath)
	return err == nil
}

// GetURL 获取本地文件访问 URL
// key: 文件存储路径
// 返回: 访问 URL
// 固定使用 /uploads 作为访问前缀，与静态文件服务保持一致
func (s *LocalStorage) GetURL(key string) string {
	// 固定访问前缀
	baseURL := "/uploads"

	// 确保 baseURL 以 / 结尾
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// URL 编码 key
	encodedKey := url.PathEscape(key)
	return baseURL + encodedKey
}

// SaveUploadedFile 保存上传的文件
// file: 上传的文件句柄
// key: 存储键
// storage: 存储服务实例
// 返回: 访问 URL 和错误
func SaveUploadedFile(file *multipart.FileHeader, key string, storage Storage) (string, error) {
	// 安全检查路径
	safeKey := strings.ReplaceAll(key, "..", "")
	if safeKey != key {
		return "", fmt.Errorf("unsafe key: %s", key)
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 获取内容类型
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传文件
	return storage.Upload(safeKey, data, contentType)
}

// GenerateFileKey 生成文件存储键
// dir: 存储目录（必须是相对路径，不含 ../），空表示直接保存在根目录
// filename: 原文件名
// 返回: 存储键（基于文件名hash，同一文件名始终相同key）
func GenerateFileKey(dir string, filename string) string {
	// 安全检查：移除 dir 中的 ../
	dir = strings.ReplaceAll(dir, "..", "")

	// 清理路径分隔符
	dir = strings.Trim(dir, "/")

	ext := filepath.Ext(filename)

	// 使用 MD5 哈希生成唯一文件名（只用filename，不含时间戳，保证同一文件始终相同key）
	hash := md5.Sum([]byte(filename))
	hashStr := fmt.Sprintf("%x", hash)[:16] // 取前16位

	// 生成 key
	var key string
	if dir != "" {
		key = fmt.Sprintf("%s/%s%s", dir, hashStr, ext)
	} else {
		// 直接保存在根目录
		key = fmt.Sprintf("%s%s", hashStr, ext)
	}
	return key
}

// GenerateFileKeyWithHash 生成文件存储键（使用内容hash）
// dir: 存储目录
// filename: 原文件名
// contentHash: 文件内容MD5 hash
// 返回: 存储键（基于内容hash，同一内容始终相同key）
func GenerateFileKeyWithHash(dir string, filename string, contentHash string) string {
	// 安全检查：移除 dir 中的 ../
	dir = strings.ReplaceAll(dir, "..", "")

	// 清理路径分隔符
	dir = strings.Trim(dir, "/")

	ext := filepath.Ext(filename)

	// 使用内容hash的前16位作为文件名
	hashStr := contentHash
	if len(hashStr) > 16 {
		hashStr = hashStr[:16]
	}

	// 生成 key
	var key string
	if dir != "" {
		key = fmt.Sprintf("%s/%s%s", dir, hashStr, ext)
	} else {
		// 直接保存在根目录
		key = fmt.Sprintf("%s%s", hashStr, ext)
	}
	return key
}
