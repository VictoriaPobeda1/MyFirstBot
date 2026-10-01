package librarycript

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type DataCoin struct {
	Data []Coins
}

type Coins struct {
	Name  string
	Quote CoinQuote
}

type CoinQuote struct {
	USD USDDetails
}

type USDDetails struct {
	Price, Percent_change_1h float64
}

func GetAllCoinFromApi() DataCoin {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil,
	)

	req.Header.Add("X-CMC_PRO_API_KEY", "582ac9dde294451b849d5acbddd61e28")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Ошибка подключения к API:", err)

	}
	defer resp.Body.Close()
	var coins DataCoin
	err = json.NewDecoder(resp.Body).Decode(&coins)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", err)
	}
	if len(coins.Data) == 0 {
		slog.Error("Пришел пустой ответ от API")
		os.Exit(1)
	}
	return coins
}
