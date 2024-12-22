package main

import (
	"log"
	"writeaid-backend/app/api/routes"
	"writeaid-backend/config"
	"writeaid-backend/util"
)

func main() {
	// 初始化配置
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 初始化数据库单例对象
	util.GetDB()
	util.GetRedisClient()
	util.GetElasticSearchClient()
	r := routes.SetupRouter()
	r.Run(config.GetServerPort())
}
