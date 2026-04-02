package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
)

func CreateReport(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var report models.Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Printf("创建举报请求解析失败: 错误: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if report.Reason == "" {
		log.Printf("创建举报失败: 原因为空, reporterID=%d", userID)
		utils.Error(w, 400, "请填写举报原因")
		return
	}

	report.ReporterID = userID
	report.Status = 0

	if err := database.DB.Create(&report).Error; err != nil {
		log.Printf("创建举报失败: reporterID=%d, 错误: %v", userID, err)
		utils.Error(w, 500, "举报失败")
		return
	}

	utils.Success(w, report)
}

func GetUserReports(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var reports []models.Report
	database.DB.Where("reporter_id = ?", userID).Order("created_at DESC").Find(&reports)

	utils.Success(w, reports)
}
