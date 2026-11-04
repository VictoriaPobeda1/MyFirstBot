package apidata

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	config "botproject/internal/config"
)

type ApiFunc interface {
	GetAllCoinFromApi() DataCoin
}

type ApiRepo struct{}

func (ap *ApiRepo) GetAllCoinFromApi() DataCoin {
	cfg := config.Load()
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", cfg.API_URL, nil,
	)

	req.Header.Add("X-CMC_PRO_API_KEY", cfg.API_KEY)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Ошибка подключения к API:", "error", err)
	}
	defer resp.Body.Close()
	var coins DataCoin
	err = json.NewDecoder(resp.Body).Decode(&coins)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", "error", err)
	}
	if len(coins.Data) == 0 {
		err := errors.New("Ошибка API")
		slog.Error("Пришел пустой ответ от API", "error", err)
	} else {
		slog.Info("Монеты получены с API")
	}
	return coins
}
