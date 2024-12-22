package config

import "github.com/spf13/viper"

func GetUserServiceAddress() string {
	return viper.GetString("user_service.address")
}

func GetNotificationServiceAddress() string {
	return viper.GetString("notification_service.address")
}

func GetKnowledgeBaseServiceAddress() string {
	return viper.GetString("knowledge_base_service.address")
}

func GetUserServicePort() string {
	return viper.GetString("user_service.port")
}

func GetNotificationServicePort() string {
	return viper.GetString("notification_service.port")
}

func GetKnowledgeBaseServicePort() string {
	return viper.GetString("knowledge_base_service.port")
}
