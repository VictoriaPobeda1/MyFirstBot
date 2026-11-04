package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	api "botproject/internal/apiData"
	coins "botproject/internal/coins"
	users "botproject/internal/users"
)

type User struct {
	Id         string
	NameCoin   string
	TimeUpdate string
}

type ServiceFunc interface {
	CheckCoinDB(coinRepo coins.CoinFunc, apiRepo api.ApiFunc)
	CheckAutoUpdate(id_user string, userRepo users.UserRepo, coinRepo coins.CoinFunc) coins.Coin
	GetAllCoin(coinRepo coins.CoinFunc) http.HandlerFunc
	RegUser(userRepo users.UserFunc) http.HandlerFunc
	SelectCoin(userRepo users.UserFunc, coinRepo coins.CoinFunc) http.HandlerFunc
	StartAuto(userRepo users.UserFunc) http.HandlerFunc
	StopAuto(userRepo users.UserFunc) http.HandlerFunc
	AutoSendlerApi(userRepo users.UserFunc) http.HandlerFunc
}

type ServiceRepo struct{}

func NewServiceRepo() *ServiceRepo {
	return &ServiceRepo{}
}

func (sr *ServiceRepo) CheckCoinDB(coinRepo coins.CoinFunc, apiRepo api.ApiFunc) {
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
func (sr *ServiceRepo) CheckAutoUpdate(id_user string, userRepo users.UserRepo, coinRepo coins.CoinFunc) coins.Coin {
	userById, _ := userRepo.SelectUserByIdDB(id_user)
	coin := coins.Coin{}
	if userById.Auto_update {
		coin, _ = coinRepo.SelectCoinByIdDB(userById.Id_coin)
	}
	return coin
}
func (sr *ServiceRepo) GetAllCoin(coinRepo coins.CoinFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, _ := coinRepo.SelectAllCoinsDB()
		dataAllCoin, err := json.Marshal(result)
		if err != nil {
			slog.Error("Ошибка вывода всех монет из БД:", err)
		}
		w.Write(dataAllCoin)
	}
}
func (sr *ServiceRepo) RegUser(userRepo users.UserFunc) http.HandlerFunc {
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
func (sr *ServiceRepo) SelectCoin(userRepo users.UserFunc, coinRepo coins.CoinFunc) http.HandlerFunc {
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
func (sr *ServiceRepo) StartAuto(userRepo users.UserFunc) http.HandlerFunc {
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
func (sr *ServiceRepo) StopAuto(userRepo users.UserFunc) http.HandlerFunc {
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
func (sr *ServiceRepo) AutoSendlerApi(userRepo users.UserFunc) http.HandlerFunc {
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
