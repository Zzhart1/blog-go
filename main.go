package main

import (
	"log/slog"
	"os"

	"blog/model"
	"blog/router"
)

func main() {
	if err := model.InitDB(); err != nil {
		slog.Error("数据库初始化失败", "err", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.SetupRouter()
	slog.Info("服务启动", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("服务启动失败", "err", err)
		os.Exit(1)
	}
}
