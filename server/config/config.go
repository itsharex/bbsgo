package config

import (
	"bbsgo/database"
	"bbsgo/models"
	"strconv"
)

func GetConfig(key string) string {
	var config models.SiteConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		return ""
	}
	return config.Value
}

func GetConfigInt(key string, defaultValue int) int {
	val := GetConfig(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func GetConfigBool(key string, defaultValue bool) bool {
	val := GetConfig(key)
	if val == "" {
		return defaultValue
	}
	return val == "true" || val == "1"
}

func SetConfig(key, value string) error {
	var config models.SiteConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		config = models.SiteConfig{
			Key:   key,
			Value: value,
		}
		return database.DB.Create(&config).Error
	}
	return database.DB.Model(&config).Update("value", value).Error
}
