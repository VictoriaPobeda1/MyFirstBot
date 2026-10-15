package users

import "time"

type User struct {
	Id_user     string
	Time_update int
	Auto_update bool
	Id_coin     int
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
	Time_send_at      time.Time
}
