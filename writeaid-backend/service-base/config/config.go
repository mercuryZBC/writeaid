package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./service-base/config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return err
	}
	return nil
}

func GetDocumentStoragePath() string {
	path := viper.GetString("document_store_path")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
	} else {
		fmt.Println("Directories created successfully")
	}
	return path
}

func GetServerPort() string {
	return ":" + viper.GetString("server.port")
}

func GetDatabaseConfig() map[string]interface{} {
	return viper.GetStringMap("database")
}

func GetRedisConfig() map[string]interface{} {
	return viper.GetStringMap("redis")
}

func GetElasticSearchAddress() string {
	return viper.GetString("elasticsearch.address")
}
