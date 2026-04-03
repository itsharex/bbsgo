package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// TencentStorage 腾讯云 COS 存储服务
type TencentStorage struct {
	config  StorageConfig // 存储配置
	client  *cos.Client   // COS 客户端
	bucket  string        // 存储桶名称
	region  string        // 区域
	domain  string        // 自定义域名
}

// newTencentStorage 创建腾讯云存储服务实例
func newTencentStorage(config StorageConfig) *TencentStorage {
	return &TencentStorage{config: config}
}

// Name 返回存储类型名称
func (s *TencentStorage) Name() string {
	return "tencent"
}

// getClient 获取 COS 客户端
// 返回: 客户端和错误
func (s *TencentStorage) getClient() (*cos.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	// 检查配置
	if s.config.Tencent.SecretId == "" || s.config.Tencent.SecretKey == "" {
		log.Printf("[tencent storage] secret credentials not configured")
		return nil, fmt.Errorf("tencent secret_id or secret_key not configured")
	}
	if s.config.Tencent.Bucket == "" || s.config.Tencent.Region == "" {
		log.Printf("[tencent storage] bucket or region not configured")
		return nil, fmt.Errorf("tencent bucket or region not configured")
	}

	// 记录配置信息
	s.bucket = s.config.Tencent.Bucket
	s.region = s.config.Tencent.Region
	s.domain = s.config.Tencent.Domain

	// 构建 COS 客户端
	// 格式: https://bucket.cos.region.myqcloud.com
	baseURL, err := cos.NewBucketURL(s.bucket, s.region, true)
	if err != nil {
		log.Printf("[tencent storage] create bucket url error: %v", err)
		return nil, fmt.Errorf("failed to create bucket url: %v", err)
	}

	httpClient := &http.Client{}
	client := cos.NewClient(&cos.BaseURL{BucketURL: baseURL}, httpClient)

	s.client = client
	return client, nil
}

// Upload 上传文件到腾讯云 COS
// key: 文件存储键
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *TencentStorage) Upload(key string, data []byte, contentType string) (string, error) {
	client, err := s.getClient()
	if err != nil {
		log.Printf("[tencent storage] get client error: %v", err)
		return "", err
	}

	// 上传文件
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}

	reader := strings.NewReader(string(data))
	_, err = client.Object.Put(context.Background(), key, reader, opt)
	if err != nil {
		log.Printf("[tencent storage] upload error, key: %s, error: %v", key, err)
		return "", fmt.Errorf("failed to upload to tencent cos: %v", err)
	}

	log.Printf("[tencent storage] upload success, key: %s", key)
	return s.GetURL(key), nil
}

// Delete 删除腾讯云 COS 文件
// key: 文件存储键
// 返回: 错误
func (s *TencentStorage) Delete(key string) error {
	client, err := s.getClient()
	if err != nil {
		log.Printf("[tencent storage] delete get client error: %v", err)
		return err
	}

	_, err = client.Object.Delete(context.Background(), key)
	if err != nil {
		log.Printf("[tencent storage] delete error, key: %s, error: %v", key, err)
		return fmt.Errorf("failed to delete from tencent cos: %v", err)
	}

	log.Printf("[tencent storage] delete success, key: %s", key)
	return nil
}

// Exists 检查腾讯云 COS 文件是否存在
// key: 文件存储键
// 返回: 是否存在
func (s *TencentStorage) Exists(key string) bool {
	client, err := s.getClient()
	if err != nil {
		return false
	}
	_, err = client.Object.Head(context.Background(), key, nil)
	return err == nil
}

// GetURL 获取腾讯云 COS 文件访问 URL
// key: 文件存储键
// 返回: 访问 URL
func (s *TencentStorage) GetURL(key string) string {
	if s.domain != "" {
		return fmt.Sprintf("https://%s/%s", s.domain, key)
	}
	if s.bucket == "" || s.region == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", s.bucket, s.region, key)
}

// TestConnection 测试腾讯云 COS 连接
// 返回: 错误
func (s *TencentStorage) TestConnection() error {
	_, err := s.getClient()
	if err != nil {
		log.Printf("[tencent storage] connection test failed: %v", err)
		return err
	}
	return nil
}

// NewTencentStorageWithCheck 创建腾讯云存储并检查配置
func NewTencentStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newTencentStorage(config)
	if err := storage.TestConnection(); err != nil {
		log.Printf("[tencent storage] connection check failed: %v", err)
		return nil, err
	}
	return storage, nil
}

// UploadFile 上传文件（通过 io.Reader）
// key: 文件存储键
// reader: 文件数据读取器
// size: 文件大小
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *TencentStorage) UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error) {
	client, err := s.getClient()
	if err != nil {
		return "", err
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   contentType,
			ContentLength: size,
		},
	}

	_, err = client.Object.Put(context.Background(), key, reader, opt)
	if err != nil {
		return "", fmt.Errorf("failed to upload to tencent cos: %v", err)
	}

	return s.GetURL(key), nil
}

// UploadURL 上传网络文件到腾讯云 COS
// key: 文件存储键
// fileURL: 网络文件 URL
// 返回: 访问 URL 和错误
func (s *TencentStorage) UploadURL(key string, fileURL string) (string, error) {
	client, err := s.getClient()
	if err != nil {
		return "", err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch url: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch url, status: %d", resp.StatusCode)
	}

	opt := &cos.ObjectPutOptions{}

	_, err = client.Object.Put(context.Background(), key, resp.Body, opt)
	if err != nil {
		return "", fmt.Errorf("failed to upload to tencent cos: %v", err)
	}

	return s.GetURL(key), nil
}
