package database

import (
	"bbsgo/models"
	"log"
)

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Forum{},
		&models.Topic{},
		&models.Post{},
		&models.Like{},
		&models.Favorite{},
		&models.Follow{},
		&models.Message{},
		&models.Notification{},
		&models.Tag{},
		&models.Report{},
		&models.Badge{},
		&models.UserBadge{},
		&models.SiteConfig{},
		&models.Draft{},
		&models.Announcement{},
		&models.VerificationCode{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")
}
