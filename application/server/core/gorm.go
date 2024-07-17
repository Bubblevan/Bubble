package core

import (
	"backend/global"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() *gorm.DB {
	// 检查是否配置了 MySQL 主机地址
	if global.Config.Mysql.Host == "" {
		global.Log.Warnln("未配置mysql，取消gorm连接")
		return nil
	}

	// 生成数据源名称（DSN）
	dsn := global.Config.Mysql.Dsn()

	// 根据环境变量设置 MySQL 日志级别
	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		// 开发环境显示所有的 SQL 语句
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		// 只打印错误的 SQL 语句
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		global.Log.Fatalf(fmt.Sprintf("[%s] mysql连接失败", dsn))
	}

	// 获取 SQL 连接
	sqlDB, _ := db.DB()
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	return db
}
