package bothandles

import (
	"botproject/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"

	models "botproject/internal/models"
)

type HandleFunc interface {
	Start(c tele.Context) error
	Rates(c tele.Context) error
	Select(c tele.Context) error
	StartAuto(c tele.Context) error
	StopAuto(c tele.Context) error
	AutoSendlerBot(b tele.Bot)
}

type HandleRepo struct {
	cfg config.Config
}

func NewHandleRepo(cfg config.Config) *HandleRepo {
	return &HandleRepo{cfg: cfg}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func (hr *HandleRepo) Start(c tele.Context) error {
	user := c.Sender()
	if user == nil {
		err := errors.New("Пустые данные")
		slog.Error("Пользователь пуст", "error", err)
	}
	data := map[string]string{"ID": fmt.Sprint(c.Sender().ID)}
	bodyBytes, _ := json.Marshal(data)
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%v/regUser", hr.cfg.APIBaseURL), bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		slog.Error("Ошибка подключения к API/regUser:", "error", err)
	}
	req.Header.Set("Content-Type", "application/json")
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Ошибка выполнения запроса к API/regUser:", "error", err)
		return err
	}
	defer result.Body.Close()
	message, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}
	return c.Send(string(message))
}
func (hr *HandleRepo) Rates(c tele.Context) error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%v/getAllCoin", hr.cfg.APIBaseURL), nil,
	)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Ошибка подключения к API/getAllCoin:", "error", err)
	}
	defer resp.Body.Close()
	var coins []models.Coins
	err = json.NewDecoder(resp.Body).Decode(&coins)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", "error", err)
	}
	args := c.Args()
	if len(args) == 0 {
		text := "Текущие курсы:\n"
		for i := 0; i < len(coins); i++ {
			name := coins[i].Name
			price := coins[i].Price
			price_min := coins[i].Price_min
			price_max := coins[i].Price_max
			percent_change_1h := coins[i].Percent_change_1h
			partText := fmt.Sprintf("\nМонета: %s\nЦена:%f$\nИзменение за час:%f%%\nМинимальная цена:%f$\nМаксимальная цена:%f$", name, price, percent_change_1h, price_min, price_max)
			lenPartText := utf8.RuneCountInString(partText)
			lenText := utf8.RuneCountInString(text)
			if lenText+lenPartText > 4096 {
				c.Send(text)
				text = partText
			} else {
				text += partText
			}
		}
	} else {
		sendCoin := strings.Join(args, " ")
		for i := range coins {
			if coins[i].Name == sendCoin {
				name := coins[i].Name
				price := coins[i].Price
				price_min := coins[i].Price_min
				price_max := coins[i].Price_max
				percent_change_1h := coins[i].Percent_change_1h
				partText := fmt.Sprintf("Текущий курс:\nМонета: %s\nЦена:%f$\nИзменение за час:%f%%\nМинимальная цена:%f$\nМаксимальная цена:%f$", name, price, percent_change_1h, price_min, price_max)
				c.Send(partText)
			}
		}
	}
	return c.Send("Все монеты выведены")
}
func (hr *HandleRepo) Select(c tele.Context) error {
	user := c.Sender()
	if user == nil {
		slog.Error("Пользователь пуст")
	}
	sendCoin := c.Args()
	if sendCoin == nil {
		slog.Error("Вы не ввели название монеты")
	}
	data := map[string]string{"NameCoin": strings.Join(sendCoin, " "), "ID": fmt.Sprint(c.Sender().ID)}
	bodyBytes, _ := json.Marshal(data)
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%v/selectCoin", hr.cfg.APIBaseURL), bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		slog.Error("Ошибка подключения к API/selectCoin:", "error", err)
	}
	req.Header.Set("Content-Type", "application/json")
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Ошибка выполнения запроса API/selectCoin:", "error", err)
	}
	defer result.Body.Close()
	message, err := io.ReadAll(result.Body)
	return c.Send(string(message))
}
func (hr *HandleRepo) StartAuto(c tele.Context) error {
	user := c.Sender()
	if user == nil {
		slog.Error("Пользователь пуст")
	}
	timeUpdate := c.Args()

	if timeUpdate == nil || !isNumeric(strings.Join(timeUpdate, "")) {
		slog.Error("Вы не задали время обновления")
		return c.Send("Не задано время обновления")

	} else {
		data := map[string]string{"TimeUpdate": fmt.Sprint(timeUpdate[0]), "ID": fmt.Sprint(c.Sender().ID)}
		bodyBytes, _ := json.Marshal(data)
		req, err := http.NewRequest(
			"GET", fmt.Sprintf("%v/startAuto", hr.cfg.APIBaseURL), bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			slog.Error("Ошибка подключения к API/startAuto:", "error", err)
		}
		req.Header.Set("Content-Type", "application/json")
		result, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Ошибка выполнения запроса API/startAuto:", "error", err)
		}
		defer result.Body.Close()
		message, err := io.ReadAll(result.Body)
		return c.Send(string(message))
	}

}
func (hr *HandleRepo) StopAuto(c tele.Context) error {
	user := c.Sender()
	if user == nil {
		slog.Error("Пользователь пуст")
	}
	data := map[string]string{"ID": fmt.Sprint(c.Sender().ID)}
	bodyBytes, _ := json.Marshal(data)
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%v/stopAuto", hr.cfg.APIBaseURL), bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		slog.Error("Ошибка подключения к API/stopAuto:", "error", err)
	}
	req.Header.Set("Content-Type", "application/json")
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Ошибка выполнения запроса API/stopAuto:", "error", err)
	}
	defer result.Body.Close()
	message, err := io.ReadAll(result.Body)
	return c.Send(string(message))
}
func (hr *HandleRepo) AutoSendlerBot(b *tele.Bot) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%v/autoSendler", hr.cfg.APIBaseURL), nil,
	)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Ошибка подключения к API/autoSendler:", "error", err)
	}
	defer resp.Body.Close()
	var usersCoin []models.UserCoin
	err = json.NewDecoder(resp.Body).Decode(&usersCoin)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", "error", err)
	}
	now := time.Now()
	for i := range usersCoin {
		timeSendAt := usersCoin[i].Time_send_at
		timeEnd := timeSendAt.Add(time.Duration(usersCoin[i].Time_update) * time.Minute)
		const layout = "2006-01-02 15:04:05"

		if now.Format(layout) >= timeEnd.Format(layout) {
			idUser, _ := strconv.Atoi(usersCoin[i].Id_user)
			name := usersCoin[i].Name
			price := usersCoin[i].Price
			percentChange1h := usersCoin[i].Percent_change_1h
			priceMin := usersCoin[i].Price_min
			priceMax := usersCoin[i].Price_max
			partText := fmt.Sprintf("Текущий курс по выбранной тобой монете такой:\nМонета: %s\nЦена:%f$\nИзменение за час:%f%%\nМинимальная цена:%f$\nМаксимальная цена:%f$", name, price, percentChange1h, priceMin, priceMax)
			_, err = b.Send(tele.ChatID(idUser), partText)
			if err != nil {
				slog.Error("Ошибка отправки:", "error", err)
			}
		}
	}
}
