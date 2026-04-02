module bbsgo

go 1.21

require (
	github.com/aliyun/aliyun-oss-go-sdk v2.1.0+incompatible
	github.com/dgraph-io/ristretto v0.1.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/mux v1.8.1
	github.com/qiniu/go-sdk/v7 v7.18.2
	github.com/tencentyun/cos-go-sdk-v5 v0.7.45
	golang.org/x/crypto v0.19.0
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
)

require (
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/glog v1.1.2 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/mozillazg/go-httpheader v0.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/time v0.15.0 // indirect
)

replace golang.org/x/time => golang.org/x/time v0.5.0

replace github.com/tencentyun/cos-go-sdk-v5 => github.com/tencentyun/cos-go-sdk-v5 v0.7.45
