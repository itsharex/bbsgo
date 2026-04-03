package antispam

import (
	"bbsgo/database"
	"bbsgo/models"
	"errors"
	"log"
	"time"
)

// ReportService 举报处理服务
// 负责处理用户举报
type ReportService struct {
	config     *ConfigService
	reputation *ReputationService
}

// NewReportService 创建举报处理服务实例
func NewReportService() *ReportService {
	return &ReportService{
		config:     GetConfigService(),
		reputation: NewReputationService(),
	}
}

// CreateReport 创建举报记录
func (s *ReportService) CreateReport(reporterID uint, targetType string, targetID uint, reason string) error {
	log.Printf("[report] creating report, reporterID: %d, targetType: %s, targetID: %d, reason: %s",
		reporterID, targetType, targetID, reason)

	maxReportsPerDay := s.config.GetInt(ConfigMaxReportsPerDay, 10)
	today := time.Now().Format("2006-01-02")

	var count int64
	if err := database.DB.Model(&models.Report{}).
		Where("reporter_id = ? AND DATE(created_at) = ?", reporterID, today).
		Count(&count).Error; err != nil {
		log.Printf("[report] query daily count failed, reporterID: %d, error: %v", reporterID, err)
		return err
	}

	if int(count) >= maxReportsPerDay {
		log.Printf("[report] daily limit reached, reporterID: %d, count: %d, limit: %d",
			reporterID, count, maxReportsPerDay)
		return errors.New("今日举报次数已达上限")
	}

	var existingReport models.Report
	result := database.DB.Where("reporter_id = ? AND target_type = ? AND target_id = ?",
		reporterID, targetType, targetID).First(&existingReport)
	if result.Error == nil {
		log.Printf("[report] duplicate report, reporterID: %d, targetType: %s, targetID: %d, existingID: %d",
			reporterID, targetType, targetID, existingReport.ID)
		return errors.New("您已举报过该内容")
	}

	report := models.Report{
		ReporterID: reporterID,
		TargetType: targetType,
		TargetID:   targetID,
		Reason:     reason,
		Status:     0,
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&report).Error; err != nil {
		log.Printf("[report] create failed, reporterID: %d, targetType: %s, targetID: %d, error: %v",
			reporterID, targetType, targetID, err)
		return err
	}

	log.Printf("[report] created, reporterID: %d, targetType: %s, targetID: %d, reportID: %d",
		reporterID, targetType, targetID, report.ID)

	go s.processReportAutoAction(targetType, targetID)

	return nil
}

// processReportAutoAction 自动处理举报
func (s *ReportService) processReportAutoAction(targetType string, targetID uint) {
	log.Printf("[report-auto] processing, targetType: %s, targetID: %d", targetType, targetID)

	threshold := s.config.GetInt(ConfigReportThreshold, 3)

	var reportCount int64
	if err := database.DB.Model(&models.Report{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&reportCount).Error; err != nil {
		log.Printf("[report-auto] query count failed, targetType: %s, targetID: %d, error: %v",
			targetType, targetID, err)
		return
	}

	log.Printf("[report-auto] count: %d, threshold: %d", reportCount, threshold)

	if int(reportCount) >= threshold {
		log.Printf("[report-auto] threshold reached, targetType: %s, targetID: %d",
			targetType, targetID)

		s.hideContent(targetType, targetID)
		s.penalizeAuthor(targetType, targetID)
	}
}

// hideContent 隐藏被举报的内容
func (s *ReportService) hideContent(targetType string, targetID uint) {
	log.Printf("[report-hide] hiding, targetType: %s, targetID: %d", targetType, targetID)

	var err error
	switch targetType {
	case "topic":
		err = database.DB.Model(&models.Topic{}).Where("id = ?", targetID).
			Update("is_hidden", true).Error
	case "comment":
		err = database.DB.Model(&models.Comment{}).Where("id = ?", targetID).
			Update("is_hidden", true).Error
	default:
		log.Printf("[report-hide] unknown type, targetType: %s, targetID: %d", targetType, targetID)
		return
	}

	if err != nil {
		log.Printf("[report-hide] failed, targetType: %s, targetID: %d, error: %v",
			targetType, targetID, err)
	} else {
		log.Printf("[report-hide] done, targetType: %s, targetID: %d", targetType, targetID)
	}
}

// penalizeAuthor 惩罚被举报内容的作者
func (s *ReportService) penalizeAuthor(targetType string, targetID uint) {
	log.Printf("[report-penalize] penalizing author, targetType: %s, targetID: %d", targetType, targetID)

	var authorID uint

	switch targetType {
	case "topic":
		var topic models.Topic
		if err := database.DB.First(&topic, targetID).Error; err != nil {
			log.Printf("[report-penalize] get topic failed, targetID: %d, error: %v", targetID, err)
			return
		}
		authorID = topic.UserID
	case "comment":
		var comment models.Comment
		if err := database.DB.First(&comment, targetID).Error; err != nil {
			log.Printf("[report-penalize] get comment failed, targetID: %d, error: %v", targetID, err)
			return
		}
		authorID = comment.UserID
	default:
		log.Printf("[report-penalize] unknown type, targetType: %s, targetID: %d", targetType, targetID)
		return
	}

	if authorID > 0 {
		if err := s.reputation.ChangeReputation(authorID, -5, "内容被举报并隐藏", targetID); err != nil {
			log.Printf("[report-penalize] deduct rep failed, authorID: %d, error: %v", authorID, err)
		} else {
			log.Printf("[report-penalize] deducted, authorID: %d, amount: 5", authorID)
		}

		go s.checkBanThreshold(authorID)
	}
}

// checkBanThreshold 检查是否需要对用户实施禁言
func (s *ReportService) checkBanThreshold(userID uint) {
	log.Printf("[report-ban-check] checking, userID: %d", userID)

	banThreshold := s.config.GetInt(ConfigReportBanThreshold, 5)
	banDays := s.config.GetInt(ConfigReportBanDays, 3)

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var hiddenTopicCount int64
	if err := database.DB.Table("topics").
		Where("user_id = ? AND is_hidden = ? AND updated_at > ?", userID, true, sevenDaysAgo).
		Count(&hiddenTopicCount).Error; err != nil {
		log.Printf("[report-ban-check] count topics failed, userID: %d, error: %v", userID, err)
	}

	var hiddenCommentCount int64
	if err := database.DB.Table("comments").
		Where("user_id = ? AND is_hidden = ? AND updated_at > ?", userID, true, sevenDaysAgo).
		Count(&hiddenCommentCount).Error; err != nil {
		log.Printf("[report-ban-check] count comments failed, userID: %d, error: %v", userID, err)
	}

	totalHidden := hiddenTopicCount + hiddenCommentCount

	log.Printf("[report-ban-check] stats, userID: %d, topics: %d, comments: %d, total: %d, threshold: %d",
		userID, hiddenTopicCount, hiddenCommentCount, totalHidden, banThreshold)

	if int(totalHidden) >= banThreshold {
		log.Printf("[report-ban-check] threshold reached, userID: %d, total: %d, threshold: %d, days: %d",
			userID, totalHidden, banThreshold, banDays)
		s.applyBan(userID, banDays)
	}
}

// applyBan 对用户实施禁言
func (s *ReportService) applyBan(userID uint, days int) {
	log.Printf("[report-ban] applying ban, userID: %d, days: %d", userID, days)

	var existingBan models.UserBan
	result := database.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Where("end_time IS NULL OR end_time > ?", time.Now()).
		First(&existingBan)

	if result.Error == nil {
		log.Printf("[report-ban] already banned, userID: %d, banID: %d", userID, existingBan.ID)
		return
	}

	endTime := time.Now().AddDate(0, 0, days)

	ban := models.UserBan{
		UserID:    userID,
		Reason:    "因多次违规被系统自动禁言",
		BanType:   "report",
		StartTime: time.Now(),
		EndTime:   &endTime,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&ban).Error; err != nil {
		log.Printf("[report-ban] create failed, userID: %d, error: %v", userID, err)
	} else {
		log.Printf("[report-ban] applied, userID: %d, banID: %d, start: %v, end: %v",
			userID, ban.ID, ban.StartTime, endTime)
	}
}

// ValidateReport 验证举报是否有效
func (s *ReportService) ValidateReport(reporterID uint, reportID uint, isValid bool) error {
	log.Printf("[report] validating, reporterID: %d, reportID: %d, isValid: %v",
		reporterID, reportID, isValid)

	var report models.Report
	if err := database.DB.First(&report, reportID).Error; err != nil {
		log.Printf("[report] report not found, reportID: %d, error: %v", reportID, err)
		return err
	}

	if report.ReporterID != reporterID {
		log.Printf("[report] unauthorized, reporterID: %d, actualReporterID: %d",
			reporterID, report.ReporterID)
		return errors.New("无权验证此举报")
	}

	if !isValid {
		log.Printf("[report] marking as malicious, reportID: %d, reporterID: %d", reportID, reporterID)
		return s.reputation.ChangeReputation(reporterID, -1, "恶意举报", reportID)
	}

	log.Printf("[report] validated, reportID: %d", reportID)
	return nil
}

// GetReportStats 获取内容的举报统计
func (s *ReportService) GetReportStats(targetType string, targetID uint) (int, error) {
	var count int64
	err := database.DB.Model(&models.Report{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error

	if err != nil {
		log.Printf("[report] get stats failed, targetType: %s, targetID: %d, error: %v",
			targetType, targetID, err)
		return 0, err
	}

	log.Printf("[report] got stats, targetType: %s, targetID: %d, count: %d",
		targetType, targetID, count)
	return int(count), nil
}

// GetUserReports 获取用户的举报列表
func (s *ReportService) GetUserReports(userID uint, page, pageSize int) ([]models.Report, int64, error) {
	log.Printf("[report] getting user reports, userID: %d, page: %d, pageSize: %d",
		userID, page, pageSize)

	var reports []models.Report
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&models.Report{}).Where("reporter_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("[report] count failed, userID: %d, error: %v", userID, err)
		return nil, 0, err
	}

	if err := database.DB.Where("reporter_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reports).Error; err != nil {
		log.Printf("[report] query failed, userID: %d, error: %v", userID, err)
		return nil, 0, err
	}

	log.Printf("[report] got reports, userID: %d, total: %d, returned: %d",
		userID, total, len(reports))
	return reports, total, nil
}
