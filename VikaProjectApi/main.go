package main

import (
	"VikaProject/logger"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	api "VikaProject/apiData"
	coins "VikaProject/coins"
	users "VikaProject/users"
	config "VikaProject/config"


	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id         string
	NameCoin   string
	TimeUpdate string
}

func CheckCoinDB(coinRepo coins.CoinFunc, apiRepo api.ApiFunc) {
	allCoins := apiRepo.GetAllCoinFromApi()
	for i := 0; i < len(allCoins.Data); i++ {
		nameCoin := allCoins.Data[i].Name
		priceCoin := allCoins.Data[i].Quote.USD.Price
		Percent_change_1h := allCoins.Data[i].Quote.USD.Percent_change_1h
		coinByName, _ := coinRepo.SelectCoinByNameDB(nameCoin)
		if coinByName.Name == "" {
			coinRepo.InsertCoinDB(nameCoin, priceCoin, Percent_change_1h, priceCoin, priceCoin)
		} else {
			coinRepo.UpdateCoinByIdDB(priceCoin, Percent_change_1h, coinByName.Id)
		}
	}
}

func CheckAutoUpdate(id_user string, userRepo users.UserRepo, coinRepo coins.CoinFunc) coins.Coin {
	userById, _ := userRepo.SelectUserByIdDB(id_user)
	coin := coins.Coin{}
	if userById.Auto_update {
		coin, _ = coinRepo.SelectCoinByIdDB(userById.Id_coin)
	}
	return coin
}

func getAllCoin(coinRepo coins.CoinFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, _ := coinRepo.SelectAllCoinsDB()
		dataAllCoin, err := json.Marshal(result)
		if err != nil {
			slog.Error("Ошибка вывода всех монет из БД:", err)
		}
		w.Write(dataAllCoin)
	}
}

func regUser(userRepo users.UserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			slog.Error("Ошибка преобразования данных:", err)
		}
		findUser, _ := userRepo.SelectUserByIdDB(user.Id)
		if findUser.Id_user == "" {
			userRepo.InsertUserDB(user.Id, 1000000, false)
			w.Write([]byte(fmt.Sprint("Спасибо за регистрацию")))
		} else {
			w.Write([]byte(fmt.Sprint("Здравствуй,UserID:" + user.Id)))
		}
	}
}

func selectCoin(userRepo users.UserFunc, coinRepo coins.CoinFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			slog.Error("Ошибка преобразования данных:", err)
		}
		findCoin, _ := coinRepo.SelectCoinByNameDB(user.NameCoin)
		if findCoin.Name == "" {
			w.Write([]byte(fmt.Sprint("Выбранной монеты нет на площадке")))
		} else {
			userRepo.UpdateCoinFromUserByIdDB(findCoin.Id, user.Id)
			w.Write([]byte(fmt.Sprint("Монета выбрана")))
		}
	}
}

func startAuto(userRepo users.UserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			slog.Error("Ошибка преобразования данных:", err)
		}
		timeUpdate, _ := strconv.Atoi(user.TimeUpdate)
		userRepo.UpdateTimeFromUserByIdDB(timeUpdate, true, user.Id)
		w.Write([]byte(fmt.Sprint("Запущено автообновление")))
	}
}

func stopAuto(userRepo users.UserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			slog.Error("Ошибка преобразования данных:", err)
		}
		userRepo.UpdateTimeFromUserByIdDB(1, false, user.Id)
		w.Write([]byte(fmt.Sprint("Выключено автообновление")))
	}
}

func autoSendler(userRepo users.UserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, _ := userRepo.SelectAllUsersByAutoUpdate()
		if len(result) > 0 {
			now := time.Now()
			for i := range result {
				timeSendAt := result[i].Time_send_at
				timeEnd := timeSendAt.Add(time.Duration(result[i].Time_update) * time.Minute)
				const layout = "2006-01-02 15:04:05"

				if now.Format(layout) >= timeEnd.Format(layout) {
					userRepo.UpdateTimeSendAtFromUserByIdDB(now, result[i].Id_user)
				}
			}
		}
		dataAllUsers, err := json.Marshal(result)
		if err != nil {
			slog.Error("Ошибка преобразования данных пользователя:", err)
		}
		w.Write(dataAllUsers)
	}
}

func main() {
	cfg := config.Load()
	if cfg.DSN == ""{
		slog.Error("Строка подключения пустая")
	}
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		slog.Error("Ошибка подключения к БД:", err)
	}
	defer pool.Close()
	userRepo := users.NewUserRepo(pool)
	coinRepo := coins.NewCoinRepo(pool)
	apiRepo := &api.ApiRepo{}

	logger.Setup("dev")
	slog.Info("Проект запущен")
	http.HandleFunc("/api/getAllCoin", getAllCoin(coinRepo))
	http.HandleFunc("/api/regUser", regUser(userRepo))
	http.HandleFunc("/api/selectCoin", selectCoin(userRepo, coinRepo))
	http.HandleFunc("/api/startAuto", startAuto(userRepo))
	http.HandleFunc("/api/stopAuto", stopAuto(userRepo))
	http.HandleFunc("/api/autoSendler", autoSendler(userRepo))
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			CheckCoinDB(coinRepo, apiRepo)
		}
	}()
	http.ListenAndServe(":8080", nil)
}
