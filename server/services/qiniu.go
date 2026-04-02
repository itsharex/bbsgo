package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"bbsgo/config"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuService struct {
	accessKey string
	secretKey string
	bucket    string
	domain    string
	mac       *qbox.Mac
	uploadMgr *storage.FormUploader
	bucketMgr *storage.BucketManager
}

func NewQiniuService() *QiniuService {
	accessKey := config.GetConfig("qiniu_access_key")
	secretKey := config.GetConfig("qiniu_secret_key")
	bucket := config.GetConfig("qiniu_bucket")
	domain := config.GetConfig("qiniu_domain")

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := &storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	return &QiniuService{
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		domain:    domain,
		mac:       mac,
		uploadMgr: storage.NewFormUploader(cfg),
		bucketMgr: storage.NewBucketManager(mac, cfg),
	}
}

func (s *QiniuService) UploadFile(key string, data io.Reader, size int64) (string, error) {
	log.Printf("[Qiniu] ========== Upload Start ==========")
	log.Printf("[Qiniu] Key: %s, Size: %d bytes", key, size)

	if s.accessKey == "" {
		log.Printf("[Qiniu] ERROR: AccessKey is empty")
		return "", fmt.Errorf("qiniu access_key not configured")
	}
	if s.secretKey == "" {
		log.Printf("[Qiniu] ERROR: SecretKey is empty")
		return "", fmt.Errorf("qiniu secret_key not configured")
	}
	if s.bucket == "" {
		log.Printf("[Qiniu] ERROR: Bucket is empty")
		return "", fmt.Errorf("qiniu bucket not configured")
	}
	if s.domain == "" {
		log.Printf("[Qiniu] ERROR: Domain is empty")
		return "", fmt.Errorf("qiniu domain not configured")
	}

	log.Printf("[Qiniu] Config: AccessKey=%s, Bucket=%s, Domain=%s", s.accessKey, s.bucket, s.domain)

	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)
	log.Printf("[Qiniu] Upload token generated: %s...", upToken[:50])

	ret := storage.PutRet{}
	err := s.uploadMgr.Put(context.Background(), &ret, upToken, key, data, size, nil)
	if err != nil {
		log.Printf("[Qiniu] ERROR: Upload failed: %v", err)
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	log.Printf("[Qiniu] Upload success! Key=%s, Hash=%s", ret.Key, ret.Hash)
	publicURL := fmt.Sprintf("https://%s/%s", s.domain, ret.Key)
	log.Printf("[Qiniu] Public URL: %s", publicURL)
	log.Printf("[Qiniu] ========== Upload End ==========")

	return publicURL, nil
}

func (s *QiniuService) UploadLocalFile(key, localPath string) (string, error) {
	log.Printf("[Qiniu] ========== Upload Local File Start ==========")
	log.Printf("[Qiniu] Key: %s, LocalPath: %s", key, localPath)

	if s.accessKey == "" {
		log.Printf("[Qiniu] ERROR: AccessKey is empty")
		return "", fmt.Errorf("qiniu access_key not configured")
	}
	if s.secretKey == "" {
		log.Printf("[Qiniu] ERROR: SecretKey is empty")
		return "", fmt.Errorf("qiniu secret_key not configured")
	}
	if s.bucket == "" {
		log.Printf("[Qiniu] ERROR: Bucket is empty")
		return "", fmt.Errorf("qiniu bucket not configured")
	}
	if s.domain == "" {
		log.Printf("[Qiniu] ERROR: Domain is empty")
		return "", fmt.Errorf("qiniu domain not configured")
	}

	log.Printf("[Qiniu] Config: AccessKey=%s, Bucket=%s, Domain=%s", s.accessKey, s.bucket, s.domain)

	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)
	log.Printf("[Qiniu] Upload token generated: %s...", upToken[:50])

	ret := storage.PutRet{}
	err := s.uploadMgr.PutFile(context.Background(), &ret, upToken, key, localPath, nil)
	if err != nil {
		log.Printf("[Qiniu] ERROR: Upload failed: %v", err)
		return "", fmt.Errorf("failed to upload local file: %v", err)
	}

	log.Printf("[Qiniu] Upload success! Key=%s, Hash=%s", ret.Key, ret.Hash)
	publicURL := fmt.Sprintf("https://%s/%s", s.domain, ret.Key)
	log.Printf("[Qiniu] Public URL: %s", publicURL)
	log.Printf("[Qiniu] ========== Upload End ==========")

	return publicURL, nil
}

func (s *QiniuService) DeleteFile(key string) error {
	err := s.bucketMgr.Delete(s.bucket, key)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

func (s *QiniuService) GetFileURL(key string) string {
	return fmt.Sprintf("https://%s/%s", s.domain, key)
}

func UploadImage(data io.Reader, size int64, filename string) (string, error) {
	service := NewQiniuService()

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("images/%d%s%s", time.Now().Unix(), fmt.Sprintf("%x", time.Now().Nanosecond()), ext)

	return service.UploadFile(key, data, size)
}

func UploadLocalImage(localPath, filename string) (string, error) {
	service := NewQiniuService()

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("images/%d%s%s", time.Now().Unix(), fmt.Sprintf("%x", time.Now().Nanosecond()), ext)

	return service.UploadLocalFile(key, localPath)
}

func UploadToQiniu(fileData []byte, fileName string, dir string) (string, error) {
	log.Printf("[Qiniu] ========== Upload Start ==========")
	log.Printf("[Qiniu] FileName: %s, Size: %d bytes, Dir: %s", fileName, len(fileData), dir)

	accessKey := config.GetConfig("qiniu_access_key")
	secretKey := config.GetConfig("qiniu_secret_key")
	bucket := config.GetConfig("qiniu_bucket")
	domain := config.GetConfig("qiniu_domain")

	log.Printf("[Qiniu] Config loaded: AccessKey=%s, Bucket=%s, Domain=%s", accessKey, bucket, domain)

	if accessKey == "" {
		log.Printf("[Qiniu] ERROR: AccessKey is empty")
		return "", fmt.Errorf("qiniu access_key not configured")
	}
	if secretKey == "" {
		log.Printf("[Qiniu] ERROR: SecretKey is empty")
		return "", fmt.Errorf("qiniu secret_key not configured")
	}
	if bucket == "" {
		log.Printf("[Qiniu] ERROR: Bucket is empty")
		return "", fmt.Errorf("qiniu bucket not configured")
	}
	if domain == "" {
		log.Printf("[Qiniu] ERROR: Domain is empty")
		return "", fmt.Errorf("qiniu domain not configured")
	}

	log.Printf("[Qiniu] Creating MAC with access key...")
	mac := qbox.NewMac(accessKey, secretKey)
	log.Printf("[Qiniu] MAC created successfully")

	log.Printf("[Qiniu] Creating put policy for bucket: %s", bucket)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	log.Printf("[Qiniu] Upload token generated: %s...", upToken[:50])

	log.Printf("[Qiniu] Creating form uploader...")
	qiniuCfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	formUploader := storage.NewFormUploader(&qiniuCfg)
	log.Printf("[Qiniu] Form uploader created")

	ext := filepath.Ext(fileName)
	if dir == "" {
		dir = "files"
	}
	key := fmt.Sprintf("%s/%d%s%s", dir, time.Now().Unix(), fmt.Sprintf("%x", time.Now().Nanosecond()), ext)
	log.Printf("[Qiniu] Upload key: %s", key)

	data := bytes.NewReader(fileData)
	log.Printf("[Qiniu] Starting upload to Qiniu...")

	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, key, data, int64(len(fileData)), nil)
	if err != nil {
		log.Printf("[Qiniu] ERROR: Upload failed: %v", err)
		return "", fmt.Errorf("upload failed: %v", err)
	}

	log.Printf("[Qiniu] Upload success! Ret: Key=%s, Hash=%s", ret.Key, ret.Hash)
	publicURL := fmt.Sprintf("https://%s/%s", domain, ret.Key)
	log.Printf("[Qiniu] Public URL: [%s]", publicURL)
	log.Printf("[Qiniu] ========== Upload End ==========")

	return publicURL, nil
}
