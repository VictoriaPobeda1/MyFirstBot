package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	API_URL       string
	API_KEY       string
	DSN           string
	TelegramToken string
	APIBaseURL    string
	TimeUpdate    time.Duration
}

func Load() Config {
	err := loadEnvFromProjectRoot()
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

func loadEnvFromProjectRoot() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		info, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil && !info.IsDir() {
			return godotenv.Load(filepath.Join(dir, ".env"))
		}
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return errors.New("корень проекта с go.mod не найден")
		}
		dir = parent
	}
}
