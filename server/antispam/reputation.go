package antispam

import (
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// ReputationService 信誉分服务
// 负责管理用户信誉分
type ReputationService struct {
	config *ConfigService
}

// NewReputationService 创建信誉分服务实例
func NewReputationService() *ReputationService {
	return &ReputationService{
		config: GetConfigService(),
	}
}

// ChangeReputation 修改用户信誉分
func (s *ReputationService) ChangeReputation(userID uint, change int, reason string, relatedID uint) error {
	log.Printf("[reputation] changing reputation, userID: %d, change: %d, reason: %s, relatedID: %d",
		userID, change, reason, relatedID)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("[reputation] user not found, userID: %d, error: %v", userID, err)
		return err
	}

	newReputation := user.Reputation + change
	if newReputation < 0 {
		newReputation = 0
		log.Printf("[reputation] reputation clipped to 0, userID: %d, old: %d, change: %d",
			userID, user.Reputation, change)
	}
	if newReputation > 100 {
		newReputation = 100
		log.Printf("[reputation] reputation clipped to 100, userID: %d, old: %d, change: %d",
			userID, user.Reputation, change)
	}

	tx := database.DB.Begin()

	if err := tx.Model(&user).Update("reputation", newReputation).Error; err != nil {
		tx.Rollback()
		log.Printf("[reputation] update failed, userID: %d, newRep: %d, error: %v",
			userID, newReputation, err)
		return err
	}

	logEntry := models.ReputationLog{
		UserID:    userID,
		Change:    change,
		Reason:    reason,
		RelatedID: relatedID,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&logEntry).Error; err != nil {
		tx.Rollback()
		log.Printf("[reputation] create log failed, userID: %d, error: %v", userID, err)
		return err
	}

	if err := s.checkAndApplyBan(tx, userID, newReputation); err != nil {
		tx.Rollback()
		log.Printf("[reputation] ban check failed, userID: %d, error: %v", userID, err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[reputation] commit failed, userID: %d, error: %v", userID, err)
		return err
	}

	log.Printf("[reputation] change success, userID: %d, old: %d, new: %d, change: %d, reason: %s",
		userID, user.Reputation, newReputation, change, reason)
	return nil
}

// checkAndApplyBan 检查并应用禁言
func (s *ReputationService) checkAndApplyBan(tx *gorm.DB, userID uint, reputation int) error {
	banLowRep := s.config.GetBool(ConfigBanLowReputation, true)
	banThreshold := s.config.GetInt(ConfigBanReputationThreshold, 20)

	if banLowRep && reputation < banThreshold {
		var existingBan models.UserBan
		result := tx.Where("user_id = ? AND is_active = ?", userID, true).First(&existingBan)
		if result.Error == nil {
			log.Printf("[reputation-ban] user already banned, userID: %d, banID: %d", userID, existingBan.ID)
			return nil
		}

		ban := models.UserBan{
			UserID:    userID,
			Reason:    "信誉分过低，系统自动禁言",
			BanType:   "reputation",
			StartTime: time.Now(),
			EndTime:   nil,
			IsActive:  true,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&ban).Error; err != nil {
			log.Printf("[reputation-ban] create ban failed, userID: %d, rep: %d, error: %v",
				userID, reputation, err)
			return err
		}

		log.Printf("[reputation-ban] user banned, userID: %d, rep: %d, threshold: %d, reason: %s",
			userID, reputation, banThreshold, ban.Reason)
	}

	return nil
}

// GetUserReputation 获取用户当前信誉分
func (s *ReputationService) GetUserReputation(userID uint) (int, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("[reputation] get failed, userID: %d, error: %v", userID, err)
		return 0, err
	}
	return user.Reputation, nil
}

// GetReputationLevel 根据信誉分获取等级名称
func (s *ReputationService) GetReputationLevel(reputation int) string {
	level := ""
	switch {
	case reputation >= 80:
		level = "normal"
	case reputation >= 60:
		level = "captcha"
	case reputation >= 40:
		level = "limited"
	case reputation >= 20:
		level = "restricted"
	default:
		level = "banned"
	}
	log.Printf("[reputation] get level, rep: %d, level: %s", reputation, level)
	return level
}

// NeedsCaptcha 判断用户是否需要验证码
func (s *ReputationService) NeedsCaptcha(userID uint) bool {
	reputation, err := s.GetUserReputation(userID)
	if err != nil {
		log.Printf("[reputation] captcha check failed, userID: %d, error: %v", userID, err)
		return false
	}

	needsCaptcha := reputation >= 60 && reputation < 80
	log.Printf("[reputation] captcha check, userID: %d, rep: %d, needsCaptcha: %v",
		userID, reputation, needsCaptcha)
	return needsCaptcha
}

// GetDailyLimits 根据信誉分获取用户每日操作限制
func (s *ReputationService) GetDailyLimits(userID uint) (topicLimit, commentLimit int) {
	reputation, err := s.GetUserReputation(userID)
	if err != nil {
		log.Printf("[reputation] daily limits failed, userID: %d, error: %v, using defaults", userID, err)
		return 10, 50
	}

	baseTopicLimit := s.config.GetInt(ConfigMaxTopicsPerDay, 10)
	baseCommentLimit := s.config.GetInt(ConfigMaxCommentsPerDay, 50)

	var topicLimitVal, commentLimitVal int
	switch {
	case reputation >= 80:
		topicLimitVal = baseTopicLimit
		commentLimitVal = baseCommentLimit
	case reputation >= 60:
		topicLimitVal = baseTopicLimit
		commentLimitVal = baseCommentLimit
	case reputation >= 40:
		topicLimitVal = 3
		commentLimitVal = 10
	default:
		topicLimitVal = 0
		commentLimitVal = 0
	}

	log.Printf("[reputation] daily limits, userID: %d, rep: %d, topicLimit: %d, commentLimit: %d",
		userID, reputation, topicLimitVal, commentLimitVal)
	return topicLimitVal, commentLimitVal
}

// AwardDailyRecovery 每日信誉恢复
func (s *ReputationService) AwardDailyRecovery() error {
	today := time.Now().Format("2006-01-02")
	log.Printf("[reputation-daily] starting daily recovery, date: %s", today)

	var users []models.User
	if err := database.DB.Where("reputation < 100").Find(&users).Error; err != nil {
		log.Printf("[reputation-daily] query users failed, error: %v", err)
		return err
	}

	recoveredCount := 0
	for _, user := range users {
		var violationCount int64
		database.DB.Model(&models.ReputationLog{}).
			Where("user_id = ? AND DATE(created_at) = ? AND change < 0", user.ID, today).
			Count(&violationCount)

		if violationCount > 0 {
			log.Printf("[reputation-daily] user has violations, skip, userID: %d, violations: %d",
				user.ID, violationCount)
			continue
		}

		var alreadyAwarded int64
		database.DB.Model(&models.ReputationLog{}).
			Where("user_id = ? AND DATE(created_at) = ? AND reason = ?", user.ID, today, "每日信誉恢复").
			Count(&alreadyAwarded)

		if alreadyAwarded > 0 {
			log.Printf("[reputation-daily] user already awarded today, skip, userID: %d", user.ID)
			continue
		}

		if err := s.ChangeReputation(user.ID, 1, "每日信誉恢复", 0); err != nil {
			log.Printf("[reputation-daily] award failed, userID: %d, error: %v", user.ID, err)
			continue
		}

		recoveredCount++
		log.Printf("[reputation-daily] awarded, userID: %d", user.ID)
	}

	log.Printf("[reputation-daily] completed, total: %d, recovered: %d", len(users), recoveredCount)
	return nil
}

// GetReputationLogs 获取用户信誉分变化日志
func (s *ReputationService) GetReputationLogs(userID uint, page, pageSize int) ([]models.ReputationLog, int64, error) {
	var logs []models.ReputationLog
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&models.ReputationLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("[reputation] query logs count failed, userID: %d, error: %v", userID, err)
		return nil, 0, err
	}

	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		log.Printf("[reputation] query logs failed, userID: %d, error: %v", userID, err)
		return nil, 0, err
	}

	log.Printf("[reputation] query logs success, userID: %d, total: %d, page: %d, pageSize: %d",
		userID, total, page, pageSize)
	return logs, total, nil
}

// RecoverFromBan 检查并恢复用户禁言状态
func (s *ReputationService) RecoverFromBan(userID uint) error {
	log.Printf("[reputation-ban-recover] checking ban status, userID: %d", userID)

	var ban models.UserBan
	result := database.DB.Where("user_id = ? AND is_active = ? AND ban_type = ?", userID, true, "reputation").
		Where("end_time IS NULL OR end_time > ?", time.Now()).
		First(&ban)

	if result.Error != nil {
		log.Printf("[reputation-ban-recover] user not banned, userID: %d", userID)
		return nil
	}

	reputation, err := s.GetUserReputation(userID)
	if err != nil {
		log.Printf("[reputation-ban-recover] get rep failed, userID: %d, error: %v", userID, err)
		return err
	}

	threshold := s.config.GetInt(ConfigBanReputationThreshold, 20)
	if reputation >= threshold {
		if err := database.DB.Model(&ban).Update("is_active", false).Error; err != nil {
			log.Printf("[reputation-ban-recover] recover failed, userID: %d, error: %v", userID, err)
			return err
		}
		log.Printf("[reputation-ban-recover] recovered, userID: %d, rep: %d, threshold: %d",
			userID, reputation, threshold)
	} else {
		log.Printf("[reputation-ban-recover] rep too low, skip, userID: %d, rep: %d, threshold: %d",
			userID, reputation, threshold)
	}

	return nil
}
