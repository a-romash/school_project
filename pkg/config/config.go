package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type IConfig struct {
	Server_address string
	Server_port    string
	Postgres_conn  string
}

var Config *IConfig
var once sync.Once

func Init() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			slog.Error("Error loading .env file")
			os.Exit(1)
		}

		Config = &IConfig{
			Server_address: os.Getenv("SERVER_ADDRESS"),
			Server_port:    os.Getenv("SERVER_PORT"),
			Postgres_conn:  os.Getenv("POSTGRES_CONN"),
		}
		if Config.Postgres_conn == "" {
			slog.Error("POSTGRES_CONN missed in .env")
			os.Exit(1)
		}
		if Config.Server_address == "" && Config.Server_port == "" {
			slog.Error("SERVER_ADDRESS and SERVER_PORT missed in .env")
			os.Exit(1)
		}
	})
}
