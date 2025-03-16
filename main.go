package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/naokij/qor5boot/admin"
)

func main() {
	// 加载 .env 文件（如果存在）
	godotenv.Load()

	h := admin.Router(admin.ConnectDB())

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
