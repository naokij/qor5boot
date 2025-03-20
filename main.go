package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/naokij/qor5boot/admin"
)

func main() {
	// 加载 .env 文件（如果存在）
	godotenv.Load()

	// 连接数据库并初始化路由
	db := admin.ConnectDB()
	h := admin.Router(db)

	// 注册信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("接收到信号: %v, 正在优雅关闭...\n", sig)

		// 关闭定时任务管理器
		admin.StopRecurringJobManager()

		fmt.Println("应用已关闭，退出")
		os.Exit(0)
	}()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	addr := host + ":" + port

	fmt.Println("Served at http://" + addr)

	mux := http.NewServeMux()
	mux.Handle("/",
		middleware.RequestID(
			middleware.Logger(
				middleware.Recoverer(h),
			),
		),
	)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
