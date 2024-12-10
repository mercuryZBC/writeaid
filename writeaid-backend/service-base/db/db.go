package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"yuqueppbackend/service-base/config"
	"yuqueppbackend/service-base/models"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// GetDB 返回单例数据库实例
func GetDB() *gorm.DB {
	// 确保初始化只执行一次
	once.Do(func() {
		// 数据库连接字符串
		config.InitConfig()
		dbConfig := config.GetDatabaseConfig()
		host := dbConfig["host"].(string)
		port := dbConfig["port"].(string)
		admin := dbConfig["user"].(string)
		password := dbConfig["password"].(string)
		dbname := dbConfig["dbname"].(string)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			admin, password, host, port, dbname)
		var err error

		// 连接数据库
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
