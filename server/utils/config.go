package utils

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"strconv"
	"time"
)

// GetConfigString 获取字符串类型的配置值
func GetConfigString(key, defaultValue string) string {
	cacheKey := "config:" + key
	if cached, ok := cache.Get(cacheKey); ok {
		if val, ok := cached.(string); ok {
			return val
		}
	}

	var config models.SiteConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err == nil {
		cache.Set(cacheKey, config.Value, 5*time.Minute)
		return config.Value
	}

	return defaultValue
}

// GetConfigBool 获取布尔类型的配置值
func GetConfigBool(key string, defaultValue bool) bool {
	strVal := GetConfigString(key, "")
	if strVal == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// GetConfigInt 获取整数类型的配置值
func GetConfigInt(key string, defaultValue int) int {
	strVal := GetConfigString(key, "")
	if strVal == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		log.Printf("get config int: failed to parse %s, error: %v", key, err)
		return defaultValue
	}
	return intVal
}

// InvalidateConfigCache 使配置缓存失效
func InvalidateConfigCache() {
	// 这里可以清空所有 config 前缀的缓存
}
