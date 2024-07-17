package utils

import (
	"backend/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "traveler:hangzhou@tcp(127.0.0.1:3306)/ASGTicket?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接到数据库失败: %v", err)
	}

	// 自动迁移
	DB.AutoMigrate(&models.User{}, &models.Ticket{}, &models.Order{})
}
