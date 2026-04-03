package config

import (
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"strconv"
	"sync"
	"time"
)

// configCache 配置缓存
var configCache = make(map[string]string)
var cacheMutex sync.RWMutex
var cacheLastLoad time.Time

// loadConfigCache 加载配置到缓存
func loadConfigCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 如果 1 分钟内已加载，直接返回
	if time.Since(cacheLastLoad) < 1*time.Minute && len(configCache) > 0 {
		return
	}

	var configs []models.SiteConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		log.Printf("load config cache failed, error: %v", err)
		return
	}

	configCache = make(map[string]string)
	for _, c := range configs {
		configCache[c.Key] = c.Value
	}
	cacheLastLoad = time.Now()
	log.Printf("[config] cache loaded, %d items", len(configCache))
}

// GetConfig 获取配置项的值
// key: 配置键名
// 返回: 配置值字符串，如果不存在则返回空字符串
func GetConfig(key string) string {
	cacheMutex.RLock()
	if val, ok := configCache[key]; ok {
		cacheMutex.RUnlock()
		return val
	}
	cacheMutex.RUnlock()

	// 缓存未命中，从数据库加载
	loadConfigCache()

	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return configCache[key]
}

// GetConfigInt 获取配置项的值作为整数
// key: 配置键名
// defaultValue: 默认值，当配置不存在或解析失败时使用
// 返回: 配置值整数
func GetConfigInt(key string, defaultValue int) int {
	val := GetConfig(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("parse config int failed, key: %s, value: %s, error: %v", key, val, err)
		return defaultValue
	}
	return i
}

// GetConfigBool 获取配置项的值作为布尔值
// key: 配置键名
// defaultValue: 默认值
// 返回: 配置值布尔值
func GetConfigBool(key string, defaultValue bool) bool {
	val := GetConfig(key)
	if val == "" {
		return defaultValue
	}
	return val == "true" || val == "1"
}

// SetConfig 设置配置项的值
// key: 配置键名
// value: 配置值
// 返回: 错误信息
func SetConfig(key, value string) error {
	var config models.SiteConfig
	// 尝试查找现有配置
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		// 配置不存在，创建新记录
		config = models.SiteConfig{
			Key:   key,
			Value: value,
		}
		if err := database.DB.Create(&config).Error; err != nil {
			log.Printf("create config failed, key: %s, value: %s, error: %v", key, value, err)
			return err
		}
	} else {
		// 配置存在，更新值
		if err := database.DB.Model(&config).Update("value", value).Error; err != nil {
			log.Printf("update config failed, key: %s, value: %s, error: %v", key, value, err)
			return err
		}
	}

	// 更新缓存
	cacheMutex.Lock()
	configCache[key] = value
	cacheMutex.Unlock()

	return nil
}

// InitConfigCache 初始化配置缓存（启动时调用）
func InitConfigCache() {
	loadConfigCache()
}
