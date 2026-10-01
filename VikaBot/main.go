package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Coins struct {
	Id                                             int
	Name                                           string
	Price, Price_min, Price_max, Percent_change_1h float64
}

type UserCoin struct {
	Id_coin           int
	Name              string
	Price             float64
	Percent_change_1h float64
	Price_min         float64
	Price_max         float64
	Auto_update       bool
	Time_update       int
	Id_user           string
}

func autoSendler(b *tele.Bot) {
	//получить данные из АПИ пользователей у которых включено тру
	//у кого тру отправлять им данные по выбранной монете

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", "http://localhost:8080/api/autoSendler", nil,
	)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка подключения к API:", err)
	}
	defer resp.Body.Close()
	var usersCoin []UserCoin
	err = json.NewDecoder(resp.Body).Decode(&usersCoin)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", err)
	}
	fmt.Println(usersCoin)
	for i := range usersCoin {
		idUser, _ := strconv.Atoi(usersCoin[i].Id_user)
		// 	idCoin := users[i].Id_coin
		// 	timeUpdate := users[i].Time_update
		// 	name := users[i].Name
		// 	price := users[i].Price
		// 	percentChange1h := users[i].Percent_change_1h
		// 	priceMin := users[i].Price_min
		// 	priceMax := users[i].Price_max
		// 	autoUpdate := users[i].Auto_update

		_, err = b.Send(tele.ChatID(idUser), "Привет! Это сообщение отправлено по твоему ID.")
		if err != nil {
			fmt.Println("Ошибка отправки: %v", err)
		}

	}
	// return c.Send("мяу")

}

func main() {
	token := "8764262646:AAF3E0oNDE5GoomxlItfWhPII9nFctsxCCQ"
	if token == "" {
		fmt.Println("Ошибка токена")
	}
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		fmt.Println("Ошибка создания бота")
		return
	}
	bot.Handle("/start", func(c tele.Context) error {
		user := c.Sender()
		if user == nil {
			fmt.Println("Пользователь пуст")
		}
		data := map[string]string{"ID": fmt.Sprint(c.Sender().ID)}
		bodyBytes, _ := json.Marshal(data)
		req, err := http.NewRequest(
			"GET", "http://localhost:8080/api/regUser", bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			fmt.Println("Ошибка подключения к API:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса API:", err)
		}
		fmt.Println(resp)
		return c.Send("Бот запущен")
	})
	bot.Handle("/rates", func(c tele.Context) error {
		client := &http.Client{}
		req, err := http.NewRequest(
			"GET", "http://localhost:8080/api/getAllCoin", nil,
		)
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка подключения к API:", err)
		}
		defer resp.Body.Close()
		var coins []Coins
		err = json.NewDecoder(resp.Body).Decode(&coins)
		if err != nil {
			slog.Error("Ошибка преобразования данных:", err)
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
		return c.Send("✅")
	})
	bot.Handle("/select", func(c tele.Context) error {
		user := c.Sender()
		if user == nil {
			fmt.Println("Пользователь пуст")
		}
		sendCoin := c.Args()
		if sendCoin == nil {
			fmt.Println("Вы не ввели название монеты")
		}
		data := map[string]string{"NameCoin": strings.Join(sendCoin, " "), "ID": fmt.Sprint(c.Sender().ID)}
		bodyBytes, _ := json.Marshal(data)
		req, err := http.NewRequest(
			"GET", "http://localhost:8080/api/selectCoin", bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			fmt.Println("Ошибка подключения к API:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса API:", err)
		}
		fmt.Println(resp)
		return c.Send("Монета выбрана")
	})
	bot.Handle("/startauto", func(c tele.Context) error {
		user := c.Sender()
		if user == nil {
			fmt.Println("Пользователь пуст")
		}
		timeUpdate := c.Args()
		if timeUpdate == nil {
			fmt.Println("Вы не задали время обновления")
		}
		data := map[string]string{"TimeUpdate": fmt.Sprint(timeUpdate[0]), "ID": fmt.Sprint(c.Sender().ID)}
		bodyBytes, _ := json.Marshal(data)
		req, err := http.NewRequest(
			"GET", "http://localhost:8080/api/startAuto", bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			fmt.Println("Ошибка подключения к API:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса API:", err)
		}
		fmt.Println(resp)
		return c.Send("Запущено автообновление")
	})
	bot.Handle("/stopauto", func(c tele.Context) error {
		user := c.Sender()
		if user == nil {
			fmt.Println("Пользователь пуст")
		}
		data := map[string]string{"ID": fmt.Sprint(c.Sender().ID)}
		bodyBytes, _ := json.Marshal(data)
		req, err := http.NewRequest(
			"GET", "http://localhost:8080/api/stopAuto", bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			fmt.Println("Ошибка подключения к API:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса API:", err)
		}
		fmt.Println(resp)
		return c.Send("Автообновление выключено")
	})
	autoSendler(bot)
	bot.Start()
}
