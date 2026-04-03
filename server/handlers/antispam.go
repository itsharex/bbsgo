package handlers

import (
	"bbsgo/antispam"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetAntiSpamConfig(w http.ResponseWriter, r *http.Request) {
	config := antispam.GetConfigService()
	allConfigs := config.GetAll()
	utils.Success(w, allConfigs)
}

func UpdateAntiSpamConfig(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	config := antispam.GetConfigService()
	for key, value := range req {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			if v {
				strValue = "true"
			} else {
				strValue = "false"
			}
		default:
			jsonBytes, _ := json.Marshal(v)
			strValue = string(jsonBytes)
		}

		if err := config.Set(key, strValue); err != nil {
			utils.Error(w, 500, "保存配置失败")
			return
		}
	}

	utils.Success(w, map[string]string{"message": "配置保存成功"})
}

func AdjustUserReputation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var req struct {
		Change int    `json:"change"`
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	reputationService := antispam.NewReputationService()
	if err := reputationService.ChangeReputation(uint(userID), req.Change, req.Reason, 0); err != nil {
		utils.Error(w, 500, "调整信誉分失败")
		return
	}

	utils.Success(w, map[string]string{"message": "信誉分调整成功"})
}

func AdminBanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var req struct {
		Reason string `json:"reason"`
		Days   int    `json:"days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Reason = "管理员手动禁言"
	}

	if req.Days == 0 {
		req.Days = 7
	}

	var existingBan models.UserBan
	result := database.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Where("end_time IS NULL OR end_time > ?", time.Now()).
		First(&existingBan)
	if result.Error == nil {
		utils.Error(w, 400, "用户已被禁言")
		return
	}

	endTime := time.Now().AddDate(0, 0, req.Days)
	ban := models.UserBan{
		UserID:    uint(userID),
		Reason:    req.Reason,
		BanType:   "admin",
		StartTime: time.Now(),
		EndTime:   &endTime,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&ban).Error; err != nil {
		utils.Error(w, 500, "禁言失败")
		return
	}

	utils.Success(w, map[string]string{"message": "用户已被禁言"})
}

func UnbanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	if err := database.DB.Model(&models.UserBan{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error; err != nil {
		utils.Error(w, 500, "解禁失败")
		return
	}

	utils.Success(w, map[string]string{"message": "用户已解禁"})
}

func GetUserReputationLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20

	reputationService := antispam.NewReputationService()
	logs, total, err := reputationService.GetReputationLogs(uint(userID), page, pageSize)
	if err != nil {
		utils.Error(w, 500, "获取信誉分日志失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":  logs,
		"total": total,
		"page":  page,
	})
}

func GetUserBanStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var ban models.UserBan
	result := database.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Where("end_time IS NULL OR end_time > ?", time.Now()).
		Order("created_at DESC").
		First(&ban)

	if result.Error != nil {
		utils.Success(w, map[string]interface{}{
			"is_banned": false,
		})
		return
	}

	utils.Success(w, map[string]interface{}{
		"is_banned": true,
		"ban":       ban,
	})
}

func GetAntiSpamStats(w http.ResponseWriter, r *http.Request) {
	var totalOperations int64
	database.DB.Model(&models.UserOperation{}).Count(&totalOperations)

	var todayOperations int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.UserOperation{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayOperations)

	var lowQualityCount int64
	database.DB.Model(&models.ContentQuality{}).
		Where("is_low_quality = ?", true).
		Count(&lowQualityCount)

	var bannedUsers int64
	database.DB.Model(&models.UserBan{}).
		Where("is_active = ?", true).
		Where("end_time IS NULL OR end_time > ?", time.Now()).
		Count(&bannedUsers)

	var totalReports int64
	database.DB.Model(&models.Report{}).Count(&totalReports)

	var pendingReports int64
	database.DB.Model(&models.Report{}).
		Where("status = ?", "pending").
		Count(&pendingReports)

	utils.Success(w, map[string]interface{}{
		"total_operations":  totalOperations,
		"today_operations":  todayOperations,
		"low_quality_count": lowQualityCount,
		"banned_users":      bannedUsers,
		"total_reports":     totalReports,
		"pending_reports":   pendingReports,
	})
}
