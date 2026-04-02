package utils

import (
	"bbsgo/config"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func GetQiniuUploadToken() string {
	accessKey := config.GetConfig("qiniu_access_key")
	secretKey := config.GetConfig("qiniu_secret_key")
	bucket := config.GetConfig("qiniu_bucket")
	
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	return putPolicy.UploadToken(mac)
}
