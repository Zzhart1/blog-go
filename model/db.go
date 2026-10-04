package model

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接与表结构，失败返回 error，由 main 统一处理
func InitDB() error {
	// 加载 .env（文件不存在时忽略，便于本地直接运行）
	_ = godotenv.Load()

	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASS", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	name := getEnv("DB_NAME", "blog_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	DB = db

	// 显式迁移全部模型，避免依赖隐式行为
	if err := db.AutoMigrate(&User{}, &Article{}, &Category{}, &Comment{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	return nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
