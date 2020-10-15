package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
	"web_app/settings"

	"github.com/spf13/viper"

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
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	// 3.初始化MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	zap.L().Debug("mysql init success...")
	// 4.初始化redis
	if err := redis.Init(); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	zap.L().Debug("redis init success...")
	// 5.注册路由
	r := routes.Setup()
	// 6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    viper.GetString("app.port"),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
