package librarycript

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Coin struct {
	Id                int
	Name              string
	Price             float64
	Percent_change_1h float64
	Price_min         float64
	Price_max         float64
}

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
}

func InsertCoinDB(name string, price float64, percent_change_1h float64, price_min float64, price_max float64) error {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pool)
	}
	defer pool.Close()
	sql := `insert into "CriptoProject"."coin" (
	name,price,percent_change_1h,price_min,price_max) 
	values ($1,$2,$3,$4,$5)`
	result, err := pool.Exec(context.Background(), sql, name, price, percent_change_1h, price_min, price_max)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return err
}

func InsertUserDB(id_user string, time_update int, auto_update bool) error {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pool)
	}
	defer pool.Close()
	sql := `insert into "CriptoProject".users(
	id_user, time_update, auto_update)
	values ($1,$2,$3)`
	result, err := pool.Exec(context.Background(), sql, id_user, time_update, auto_update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return err
}

func SelectAllCoinsDB() ([]Coin, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(pool)
	}
	defer pool.Close()
	sql := `select id, name, price, percent_change_1h, price_min, price_max
	from "CriptoProject".coin;`
	result, err := pool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	coins := []Coin{}
	for result.Next() {
		c := Coin{}
		err := result.Scan(&c.Id, &c.Name, &c.Price, &c.Percent_change_1h, &c.Price_min, &c.Price_max)
		if err != nil {
			fmt.Println(err)
			continue
		}
		coins = append(coins, c)
	}
	return coins, err
}

func SelectAllUsersDB() ([]User, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(pool)
	}
	defer pool.Close()
	sql := `select id_user, time_update, auto_update, id_coin
	from "CriptoProject".users`
	result, err := pool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	users := []User{}
	for result.Next() {
		u := User{}
		err := result.Scan(&u.Id_user, &u.Time_update, &u.Auto_update, &u.Id_coin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}
	return users, err
}

func SelectCoinByIdDB(id int) (Coin, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(pool)
	}
	defer pool.Close()
	sql := `select id, name, price, percent_change_1h, price_min, price_max
	from "CriptoProject".coin where id=$1`
	coin := Coin{}
	err = pool.QueryRow(context.Background(), sql, id).Scan(&coin.Id, &coin.Name, &coin.Percent_change_1h, &coin.Price, &coin.Price_max, &coin.Price_min)
	if err != nil {
		fmt.Println(err)
	}
	return coin, err
}

func SelectCoinByNameDB(name string) (Coin, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(pool)
	}
	defer pool.Close()
	sql := `select id,name,price,percent_change_1h, price_min,price_max
	from "CriptoProject".coin where name=$1`
	coin := Coin{}
	err = pool.QueryRow(context.Background(), sql, name).Scan(&coin.Id, &coin.Name, &coin.Percent_change_1h, &coin.Price, &coin.Price_max, &coin.Price_min)
	if err != nil {
		fmt.Println(err)
	}
	return coin, err
}

func SelectUserByIdDB(id_user string) (User, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(pool)
	}
	defer pool.Close()
	sql := `select id_user, time_update, auto_update, id_coin
	from "CriptoProject".users where id_user=$1`
	user := User{}
	err = pool.QueryRow(context.Background(), sql, id_user).Scan(&user.Id_user, &user.Time_update, &user.Auto_update, &user.Id_coin)
	if err != nil {
		fmt.Println(err)
	}
	slog.Info("Пользователь выведен")
	return user, err
}

func SelectAllUsersByAutoUpdate() ([]UserCoin, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer pool.Close()
	sql := `select "CriptoProject".coin.*,"CriptoProject".users.auto_update,
	"CriptoProject".users.time_update,"CriptoProject".users.id_user 
	from "CriptoProject".users inner join "CriptoProject".coin
	on "CriptoProject".users.id_coin="CriptoProject".coin.id
	where "CriptoProject".users.auto_update=true`
	result, err := pool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	userCoin := []UserCoin{}
	for result.Next() {
		u := UserCoin{}
		err := result.Scan(&u.Id_coin, &u.Name, &u.Price, &u.Percent_change_1h, &u.Price_min, &u.Price_max, &u.Auto_update, &u.Time_update, &u.Id_user)
		if err != nil {
			fmt.Println(err)
			continue
		}
		userCoin = append(userCoin, u)
	}
	slog.Info("Пользователь выведен")
	return userCoin, err
}

func UpdateCoinFromUserByIdDB(id_coin int, id_user string) error {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pool)
	}
	defer pool.Close()
	sql := `update "CriptoProject".users
	set id_coin=$1
	where id_user=$2`
	result, err := pool.Exec(context.Background(), sql, id_coin, id_user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return err
}

func UpdateTimeFromUserByIdDB(time_update int, auto_update bool, id_user string) error {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pool)
	}
	defer pool.Close()
	sql := `update "CriptoProject".users
	set time_update=$1, auto_update=$2
	where id_user=$3`
	result, err := pool.Exec(context.Background(), sql, time_update, auto_update, id_user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return err
}

func UpdateCoinByIdDB(price float64, percent_change_1h float64, id int) (pgconn.CommandTag, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(pool)
	}
	defer pool.Close()
	coin, err := SelectCoinByIdDB(id)
	if price < coin.Price_min {
		coin.Price_min = price
	}
	if price > coin.Price_max {
		coin.Price_max = price
	}
	sql := `UPDATE "CriptoProject".coin
	SET price=$1, percent_change_1h=$2, price_min=$3, price_max=$4
	WHERE id=$5`
	result, err := pool.Exec(context.Background(), sql, price, percent_change_1h, coin.Price_min, coin.Price_max, id)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}
