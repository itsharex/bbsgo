package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunStorage 阿里云 OSS 存储服务
type AliyunStorage struct {
	config     StorageConfig // 存储配置
	client     *oss.Client   // OSS 客户端
	bucket     *oss.Bucket  // 存储桶
}

// newAliyunStorage 创建阿里云存储服务实例
func newAliyunStorage(config StorageConfig) *AliyunStorage {
	return &AliyunStorage{config: config}
}

// Name 返回存储类型名称
func (s *AliyunStorage) Name() string {
	return "aliyun"
}

// getClient 获取 OSS 客户端
// 返回: 客户端和错误
func (s *AliyunStorage) getClient() (*oss.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	// 检查配置
	if s.config.Aliyun.AccessKeyId == "" || s.config.Aliyun.AccessKeySecret == "" {
		log.Printf("[aliyun storage] access key not configured")
		return nil, fmt.Errorf("aliyun access_key not configured")
	}
	if s.config.Aliyun.Endpoint == "" {
		log.Printf("[aliyun storage] endpoint not configured")
		return nil, fmt.Errorf("aliyun endpoint not configured")
	}

	// 创建 OSS 客户端
	client, err := oss.New(s.config.Aliyun.Endpoint, s.config.Aliyun.AccessKeyId, s.config.Aliyun.AccessKeySecret)
	if err != nil {
		log.Printf("[aliyun storage] create client error: %v", err)
		return nil, fmt.Errorf("failed to create aliyun oss client: %v", err)
	}

	s.client = client
	return client, nil
}

// getBucket 获取存储桶
// 返回: 存储桶和错误
func (s *AliyunStorage) getBucket() (*oss.Bucket, error) {
	if s.bucket != nil {
		return s.bucket, nil
	}

	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	if s.config.Aliyun.Bucket == "" {
		log.Printf("[aliyun storage] bucket not configured")
		return nil, fmt.Errorf("aliyun bucket not configured")
	}

	bucket, err := client.Bucket(s.config.Aliyun.Bucket)
	if err != nil {
		log.Printf("[aliyun storage] get bucket error: %v", err)
		return nil, fmt.Errorf("failed to get aliyun oss bucket: %v", err)
	}

	s.bucket = bucket
	return bucket, nil
}

// Upload 上传文件到阿里云 OSS
// key: 文件存储键
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *AliyunStorage) Upload(key string, data []byte, contentType string) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		log.Printf("[aliyun storage] get bucket error: %v", err)
		return "", err
	}

	// 上传文件
	options := []oss.Option{
		oss.ContentType(contentType),
	}
	reader := bytes.NewReader(data)
	err = bucket.PutObject(key, reader, options...)
	if err != nil {
		log.Printf("[aliyun storage] upload error, key: %s, error: %v", key, err)
		return "", fmt.Errorf("failed to upload to aliyun oss: %v", err)
	}

	log.Printf("[aliyun storage] upload success, key: %s", key)
	return s.GetURL(key), nil
}

// Delete 删除阿里云 OSS 文件
// key: 文件存储键
// 返回: 错误
func (s *AliyunStorage) Delete(key string) error {
	bucket, err := s.getBucket()
	if err != nil {
		log.Printf("[aliyun storage] delete get bucket error: %v", err)
		return err
	}

	err = bucket.DeleteObject(key)
	if err != nil {
		log.Printf("[aliyun storage] delete error, key: %s, error: %v", key, err)
		return fmt.Errorf("failed to delete from aliyun oss: %v", err)
	}

	log.Printf("[aliyun storage] delete success, key: %s", key)
	return nil
}

// GetURL 获取阿里云 OSS 文件访问 URL
// key: 文件存储键
// 返回: 访问 URL
func (s *AliyunStorage) GetURL(key string) string {
	if s.config.Aliyun.Domain == "" {
		// 如果没有配置域名，使用默认格式
		if s.config.Aliyun.Bucket == "" || s.config.Aliyun.Endpoint == "" {
			return ""
		}
		// 处理 endpoint 格式，转换为 URL
		endpoint := s.config.Aliyun.Endpoint
		if strings.HasPrefix(endpoint, "http://") {
			endpoint = strings.TrimPrefix(endpoint, "http://")
		} else if strings.HasPrefix(endpoint, "https://") {
			endpoint = strings.TrimPrefix(endpoint, "https://")
		}
		return fmt.Sprintf("https://%s.%s/%s", s.config.Aliyun.Bucket, endpoint, key)
	}
	return fmt.Sprintf("https://%s/%s", s.config.Aliyun.Domain, key)
}

// TestConnection 测试阿里云 OSS 连接
// 返回: 错误
func (s *AliyunStorage) TestConnection() error {
	_, err := s.getBucket()
	if err != nil {
		log.Printf("[aliyun storage] connection test failed: %v", err)
		return err
	}
	return nil
}

// NewAliyunStorageWithCheck 创建阿里云存储并检查配置
func NewAliyunStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newAliyunStorage(config)
	if err := storage.TestConnection(); err != nil {
		log.Printf("[aliyun storage] connection check failed: %v", err)
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
func (s *AliyunStorage) UploadFile(key string, reader io.Reader, size int64, contentType string) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}

	options := []oss.Option{
		oss.ContentType(contentType),
		oss.ContentLength(size),
	}

	err = bucket.PutObject(key, reader, options...)
	if err != nil {
		return "", fmt.Errorf("failed to upload to aliyun oss: %v", err)
	}

	return s.GetURL(key), nil
}

// UploadURL 上传网络文件到阿里云 OSS
// key: 文件存储键
// url: 网络文件 URL
// 返回: 访问 URL 和错误
func (s *AliyunStorage) UploadURL(key string, url string) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch url: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch url, status: %d", resp.StatusCode)
	}

	err = bucket.PutObject(key, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to upload to aliyun oss: %v", err)
	}

	return s.GetURL(key), nil
}
