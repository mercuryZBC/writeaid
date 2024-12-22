package util

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
	"writeaid-backend/config"
	"writeaid-backend/models"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

// GetDB 返回单例数据库实例
func GetDB() *gorm.DB {
	// 确保初始化只执行一次
	dbOnce.Do(func() {
		// 数据库连接字符串
		dbConfig := config.GetDatabaseConfig()
		host := dbConfig["host"].(string)
		port := dbConfig["port"].(string)
		admin := dbConfig["user"].(string)
		password := dbConfig["password"].(string)
		dbname := dbConfig["dbname"].(string)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			admin, password, host, port, dbname)
		var err error
		dbLogFile, err := os.OpenFile(config.GetDBLogfilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		dbLogger := logger.New(
			log.New(dbLogFile, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)

		// 连接数据库
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		log.Println("Successfully connected to database")
		if err := models.MigrateDB(dbInstance); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
		//m := gormigrate.New(dbInstance, gormigrate.DefaultOptions, []*gormigrate.Migration{})
	})
	return dbInstance
}
