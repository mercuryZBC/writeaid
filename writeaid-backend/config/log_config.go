package config

import "github.com/spf13/viper"

func GetDBLogfilePath() string {
	addr := viper.GetString("logfile_path.db")
	return addr
}

func GetUserServiceLogfilePath() string {
	addr := viper.GetString("logfile_path.user_service")
	return addr
}
func GetNotificationServiceLogfilePath() string {
	return viper.GetString("logfile_path.notification_service")
}
func GetKnowledgeBaseServiceLogfilePath() string {
	return viper.GetString("logfile_path.knowledge_base_service")
}
