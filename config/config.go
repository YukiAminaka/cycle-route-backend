package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server Server
}

type DBConfig struct {
	Name string `env:"DB_NAME"`
	User string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
    Host string `env:"DB_HOST"`
    Port string `env:"DB_PORT"`
}

type Server struct {
	Address string `env:"ADDRESS" envDefault:"0.0.0.0"`
	Port    string `env:"PORT" envDefault:"8080"`
}

// 読み込み
var (
	cfg Config
	once sync.Once
)


func GetConfig() *Config {
	// goroutine実行中でも一度だけ実行される
	once.Do(func() {
		// envファイルの読み込み
		err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("GO_ENV")))
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		if err := env.Parse(&cfg); err != nil {
			panic(err)
		}
	})
	return &cfg
}