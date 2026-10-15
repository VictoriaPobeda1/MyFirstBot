package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	APIBaseURL    string
	TimeUpdate    time.Duration
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Ошибка загрузки файла .env")
	}
	var cfg Config
	cfg.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if cfg.TelegramToken == "" {
		slog.Error("Токен пустой")
	}
	cfg.APIBaseURL = os.Getenv("API_BASE_URL")
	if cfg.APIBaseURL == "" {
		slog.Error("Ссылка на API пустая")
	}
	timeUpdate, err := strconv.Atoi(os.Getenv("TIME_UPDATE"))
	if err != nil {
		slog.Error("Время задано неккоректно")
	}
	cfg.TimeUpdate = time.Duration(timeUpdate) * time.Minute

	return cfg
}
