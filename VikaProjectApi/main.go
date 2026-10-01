package main

import (
	"VikaProject/librarycript"
	"VikaProject/logger"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type User struct {
	Id         string
	NameCoin   string
	TimeUpdate string
}

func CheckCoinDB() {
	allCoins := librarycript.GetAllCoinFromApi()
	for i := 0; i < len(allCoins.Data); i++ {
		nameCoin := allCoins.Data[i].Name
		priceCoin := allCoins.Data[i].Quote.USD.Price
		Percent_change_1h := allCoins.Data[i].Quote.USD.Percent_change_1h
		coinByName, _ := librarycript.SelectCoinByNameDB(nameCoin)
		if coinByName.Name == "" {
			librarycript.InsertCoinDB(nameCoin, priceCoin, Percent_change_1h, priceCoin, priceCoin)
		} else {
			librarycript.UpdateCoinByIdDB(priceCoin, Percent_change_1h, coinByName.Id)
		}
	}
}

func CheckAutoUpdate(id_user string) librarycript.Coin {
	userById, _ := librarycript.SelectUserByIdDB(id_user)
	coin := librarycript.Coin{}
	if userById.Auto_update {
		coin, _ = librarycript.SelectCoinByIdDB(userById.Id_coin)
	}
	return coin
}

func getAllCoin(w http.ResponseWriter, r *http.Request) {
	result, _ := librarycript.SelectAllCoinsDB()
	dataAllCoin, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Ошибка преобразования")
	}
	w.Write(dataAllCoin)
}

func regUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", err)
	}
	findUser, _ := librarycript.SelectUserByIdDB(user.Id)
	if findUser.Id_user == "" {
		librarycript.InsertUserDB(user.Id, 1000000, false)
	} else {
		w.Write([]byte(fmt.Sprint("Здравствуй,UserID:" + user.Id)))
	}
}

func selectCoin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", err)
	}
	findCoin, _ := librarycript.SelectCoinByNameDB(user.NameCoin)
	if findCoin.Name == "" {
		w.Write([]byte(fmt.Sprint("Выбранной монеты нет на площадке")))
	} else {
		librarycript.UpdateCoinFromUserByIdDB(findCoin.Id, user.Id)
		w.Write([]byte(fmt.Sprint("Монета выбрана")))
	}
}

func startAuto(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", err)
	}
	timeUpdate, _ := strconv.Atoi(user.TimeUpdate)
	librarycript.UpdateTimeFromUserByIdDB(timeUpdate, true, user.Id)
	w.Write([]byte(fmt.Sprint("Запущено автообновление")))
}

func stopAuto(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("Ошибка преобразования данных:", err)
	}
	librarycript.UpdateTimeFromUserByIdDB(1, false, user.Id)
	w.Write([]byte(fmt.Sprint("Выключено автообновление")))
}

func autoSendler(w http.ResponseWriter, r *http.Request) {
	result, _ := librarycript.SelectAllUsersByAutoUpdate()
	dataAllUsers, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Ошибка преобразования")
	}
	w.Write(dataAllUsers)

}

func main() {
	logger.Setup("dev")
	slog.Info("Проект запущен")
	http.HandleFunc("/api/getAllCoin", getAllCoin)
	http.HandleFunc("/api/regUser", regUser)
	http.HandleFunc("/api/selectCoin", selectCoin)
	http.HandleFunc("/api/startAuto", startAuto)
	http.HandleFunc("/api/stopAuto", stopAuto)
	http.HandleFunc("/api/autoSendler", autoSendler)
	http.ListenAndServe(":8080", nil)
}
