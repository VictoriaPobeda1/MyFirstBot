package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	API_URL string
	API_KEY string
	DSN     string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Ошибка загрузки файла .env")
	}
	var cfg Config
	cfg.API_URL = os.Getenv("API_URL")
	if cfg.API_URL == "" {
		slog.Error("Ссылка на API пустая")
	}
	cfg.API_KEY = os.Getenv("API_KEY")
	if cfg.API_KEY == "" {
		slog.Error("Ключ API пустой")
	}
	cfg.DSN = os.Getenv("DSN")
	if cfg.DSN == "" {
		slog.Error("Строка подключения пустая")
	}
	return cfg
}
