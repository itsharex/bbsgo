package antispam

import (
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"math"
	"time"
)

// HotScoreService 热度评分服务
// 负责计算和管理话题的热度评分
type HotScoreService struct {
	config  *ConfigService
	quality *ContentQualityService
}

// NewHotScoreService 创建热度评分服务实例
func NewHotScoreService() *HotScoreService {
	return &HotScoreService{
		config:  GetConfigService(),
		quality: NewContentQualityService(),
	}
}

// CalculateHotScore 计算话题热度评分
func (s *HotScoreService) CalculateHotScore(topic *models.Topic) float64 {
	log.Printf("[hotscore] calculating, topicID: %d, likes: %d, comments: %d, views: %d",
		topic.ID, topic.LikeCount, topic.ReplyCount, topic.ViewCount)

	baseScore := s.calculateBaseScore(topic)
	log.Printf("[hotscore] base score, topicID: %d, baseScore: %.2f", topic.ID, baseScore)

	multiplier := 1.0

	if s.quality.IsLowQuality(topic.ID, "topic") {
		lowQualityMultiplier := s.config.GetFloat(ConfigLowQualityHotMultiplier, 0.3)
		multiplier *= lowQualityMultiplier
		log.Printf("[hotscore] low quality, topicID: %d, multiplier: %.2f", topic.ID, lowQualityMultiplier)
	}

	var user models.User
	if err := database.DB.First(&user, topic.UserID).Error; err == nil {
		lowRepThreshold := s.config.GetInt(ConfigLowReputationThreshold, 60)
		if user.Reputation < lowRepThreshold {
			lowRepMultiplier := s.config.GetFloat(ConfigLowReputationHotMultiplier, 0.5)
			multiplier *= lowRepMultiplier
			log.Printf("[hotscore] low reputation, topicID: %d, userID: %d, rep: %d, multiplier: %.2f",
				topic.ID, topic.UserID, user.Reputation, lowRepMultiplier)
		}
	} else {
		log.Printf("[hotscore] get user failed, topicID: %d, userID: %d, error: %v",
			topic.ID, topic.UserID, err)
	}

	finalScore := baseScore * multiplier
	log.Printf("[hotscore] calculated, topicID: %d, finalScore: %.2f (base: %.2f * mult: %.2f)",
		topic.ID, finalScore, baseScore, multiplier)

	return finalScore
}

// calculateBaseScore 计算话题基础得分
func (s *HotScoreService) calculateBaseScore(topic *models.Topic) float64 {
	likeScore := float64(topic.LikeCount) * 2.0
	commentScore := float64(topic.ReplyCount) * 1.5
	viewScore := float64(topic.ViewCount) * 0.1

	totalScore := likeScore + commentScore + viewScore

	hours := time.Since(topic.CreatedAt).Hours()
	if hours < 1 {
		hours = 1
	}

	timeDecay := math.Pow(hours, 0.5)

	baseScore := totalScore / timeDecay

	log.Printf("[hotscore-base] topicID: %d, likeScore: %.2f, commentScore: %.2f, viewScore: %.2f, total: %.2f, hours: %.2f, decay: %.2f, base: %.2f",
		topic.ID, likeScore, commentScore, viewScore, totalScore, hours, timeDecay, baseScore)

	return baseScore
}

// UpdateTopicHotScores 更新所有话题的热度评分
func (s *HotScoreService) UpdateTopicHotScores() error {
	log.Printf("[hotscore] updating all topic scores")

	var topics []models.Topic
	if err := database.DB.Find(&topics).Error; err != nil {
		log.Printf("[hotscore] query topics failed, error: %v", err)
		return err
	}

	updatedCount := 0
	for _, topic := range topics {
		score := s.CalculateHotScore(&topic)
		if err := database.DB.Model(&topic).Update("hot_score", score).Error; err != nil {
			log.Printf("[hotscore] update failed, topicID: %d, error: %v", topic.ID, err)
			continue
		}
		updatedCount++
	}

	log.Printf("[hotscore] update completed, total: %d, updated: %d", len(topics), updatedCount)
	return nil
}

// GetHotTopics 获取热门话题列表
func (s *HotScoreService) GetHotTopics(forumID uint, page, pageSize int) ([]models.Topic, int64, error) {
	log.Printf("[hotscore] querying hot topics, forumID: %d, page: %d, pageSize: %d",
		forumID, page, pageSize)

	var topics []models.Topic
	var total int64

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Topic{}).Where("is_hidden = ? OR is_hidden IS NULL", false)

	if forumID > 0 {
		query = query.Where("forum_id = ?", forumID)
		log.Printf("[hotscore] filtering by forum, forumID: %d", forumID)
	}

	if err := query.Count(&total).Error; err != nil {
		log.Printf("[hotscore] count failed, error: %v", err)
		return nil, 0, err
	}

	if err := query.Order("hot_score DESC, created_at DESC").
		Preload("User").
		Preload("Forum").
		Offset(offset).
		Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("[hotscore] query failed, error: %v", err)
		return nil, 0, err
	}

	log.Printf("[hotscore] query success, forumID: %d, total: %d, page: %d, pageSize: %d, returned: %d",
		forumID, total, page, pageSize, len(topics))

	return topics, total, nil
}

// RecalculateAllScores 重新计算所有话题的热度评分
func (s *HotScoreService) RecalculateAllScores() error {
	log.Printf("[hotscore] recalculating all scores")
	return s.UpdateTopicHotScores()
}
