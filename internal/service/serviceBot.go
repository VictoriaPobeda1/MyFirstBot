package service

import (
	"botproject/internal/config"
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v3"

	handles "botproject/internal/botHandles"
)

type ServiceBotFunc interface {
	RunBot()
}

type ServiceBotRepo struct {
	cfg config.Config
}

func NewServiceBotRepo(cfg config.Config) *ServiceBotRepo {
	return &ServiceBotRepo{cfg: cfg}
}

func (sbr *ServiceBotRepo) RunBot() {
	pref := tele.Settings{
		Token:  sbr.cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("Ошибка создания бота")
		return
	}
	handleRepo := handles.NewHandleRepo(sbr.cfg)
	bot.Handle("/start", handleRepo.Start)
	bot.Handle("/rates", handleRepo.Rates)
	bot.Handle("/select", handleRepo.Select)
	bot.Handle("/startauto", handleRepo.StartAuto)
	bot.Handle("/stopauto", handleRepo.StopAuto)
	ticker := time.NewTicker(sbr.cfg.TimeUpdate)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			handleRepo.AutoSendlerBot(bot)
		}
	}()
	slog.Info("Бот запущен")
	bot.Start()
}
