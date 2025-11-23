package config

import (
	"os"
)

// Config はアプリケーションの設定
type Config struct {
	DatabaseURL string
	ServerPort  string
}

// LoadConfig は環境変数から設定を読み込む
func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
