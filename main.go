package main

import (
	"fmt"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
	"web_app/settings"

	"go.uber.org/zap"
)

// Go Web 开发通用脚手架模板
func main() {
	// 1.加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("Init Setting failed, err:%v\n", err)
		return
	}
	// 2.初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Init Logger failed, err:%v\n", err)
		return
	}
	zap.L().Debug("logger init success...")
	// 3.初始化MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
		return
	}
	// 4.初始化redis
	if err := redis.Init(); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
		return
	}
	// 5.注册路由
	r := routes.Setup()
	// 6.启动服务（优雅关机）

}
