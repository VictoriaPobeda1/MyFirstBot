package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	api "botproject/internal/apiData"
	coins "botproject/internal/coins"
	config "botproject/internal/config"
	logger "botproject/internal/logger"
	service "botproject/internal/service"
	users "botproject/internal/users"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	serviceRepo := service.NewServiceRepo()
	cfg := config.Load()
	if cfg.DSN == "" {
		slog.Error("Строка подключения пустая")
	}
	if cfg.TelegramToken == "" {
		slog.Error("Ошибка токена")
	}
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		slog.Error("Ошибка подключения к БД:", err)
	}
	defer pool.Close()
	userRepo := users.NewUserRepo(pool)
	coinRepo := coins.NewCoinRepo(pool)
	apiRepo := &api.ApiRepo{}
	serviceBot := service.NewServiceBotRepo(cfg)

	logger.Setup("dev")
	slog.Info("Проект запущен")
	http.HandleFunc("/api/getAllCoin", serviceRepo.GetAllCoin(coinRepo))
	http.HandleFunc("/api/regUser", serviceRepo.RegUser(userRepo))
	http.HandleFunc("/api/selectCoin", serviceRepo.SelectCoin(userRepo, coinRepo))
	http.HandleFunc("/api/startAuto", serviceRepo.StartAuto(userRepo))
	http.HandleFunc("/api/stopAuto", serviceRepo.StopAuto(userRepo))
	http.HandleFunc("/api/autoSendler", serviceRepo.AutoSendlerApi(userRepo))
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			serviceRepo.CheckCoinDB(coinRepo, apiRepo)
		}
	}()

	go serviceBot.RunBot()
	http.ListenAndServe(":8080", nil)
}
