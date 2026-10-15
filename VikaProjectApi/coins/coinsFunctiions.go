package coins

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinFunc interface {
	InsertCoinDB(name string, price float64, percent_change_1h float64, price_min float64, price_max float64) error
	SelectAllCoinsDB() ([]Coin, error)
	SelectCoinByIdDB(id int) (Coin, error)
	SelectCoinByNameDB(name string) (Coin, error)
	UpdateCoinByIdDB(price float64, percent_change_1h float64, id int) (pgconn.CommandTag, error)
}

type CoinRepo struct {
	pool *pgxpool.Pool
}

func NewCoinRepo(pool *pgxpool.Pool) *CoinRepo {
	return &CoinRepo{pool: pool}
}

func (cr *CoinRepo) InsertCoinDB(name string, price float64, percent_change_1h float64, price_min float64, price_max float64) error {
	sql := `insert into "CriptoProject"."coin" (
	name,price,percent_change_1h,price_min,price_max) 
	values ($1,$2,$3,$4,$5)`
	_, err := cr.pool.Exec(context.Background(), sql, name, price, percent_change_1h, price_min, price_max)
	if err != nil {
		slog.Error("Ошибка добавления монеты в БД:", err)
	}
	slog.Info("Монета добавлена")
	return err
}

func (cr *CoinRepo) SelectAllCoinsDB() ([]Coin, error) {
	sql := `select id, name, price, percent_change_1h, price_min, price_max
	from "CriptoProject".coin;`
	result, err := cr.pool.Query(context.Background(), sql)
	if err != nil {
		slog.Error("Ошибка вывода монет из БД:", err)
	}
	defer result.Close()
	coins := []Coin{}
	for result.Next() {
		c := Coin{}
		err := result.Scan(&c.Id, &c.Name, &c.Price, &c.Percent_change_1h, &c.Price_min, &c.Price_max)
		if err != nil {
			slog.Error("Ошибка считывания данных о монете:", err)
			continue
		}
		coins = append(coins, c)
	}
	slog.Info("Данные о монетах выведены из БД")
	return coins, err
}

func (cr *CoinRepo) SelectCoinByIdDB(id int) (Coin, error) {
	sql := `select id, name, price, percent_change_1h, price_min, price_max
	from "CriptoProject".coin where id=$1`
	coin := Coin{}
	err := cr.pool.QueryRow(context.Background(), sql, id).Scan(&coin.Id, &coin.Name, &coin.Percent_change_1h, &coin.Price, &coin.Price_max, &coin.Price_min)
	if err != nil {
		slog.Error("Ошибка вывода монеты из БД по id:", err)
	}
	slog.Info("Данные о монете выведены из БД по id")
	return coin, err
}

func (cr *CoinRepo) SelectCoinByNameDB(name string) (Coin, error) {
	sql := `select id,name,price,percent_change_1h, price_min,price_max
	from "CriptoProject".coin where name=$1`
	coin := Coin{}
	err := cr.pool.QueryRow(context.Background(), sql, name).Scan(&coin.Id, &coin.Name, &coin.Percent_change_1h, &coin.Price, &coin.Price_max, &coin.Price_min)
	if err != nil {
		slog.Error("Ошибка вывода монеты из БД по имени:", err)
	}
	slog.Info("Данные о монете выведены из БД по имени")
	return coin, err
}

func (cr *CoinRepo) UpdateCoinByIdDB(price float64, percent_change_1h float64, id int) (pgconn.CommandTag, error) {
	coin, err := cr.SelectCoinByIdDB(id)
	if price < coin.Price_min {
		coin.Price_min = price
	}
	if price > coin.Price_max {
		coin.Price_max = price
	}
	sql := `update "CriptoProject".coin
	set price=$1, percent_change_1h=$2, price_min=$3, price_max=$4
	where id=$5`
	result, err := cr.pool.Exec(context.Background(), sql, price, percent_change_1h, coin.Price_min, coin.Price_max, id)
	if err != nil {
		slog.Error("Ошибка обновления данных о монете в БД:", err)
	}
	slog.Info("Обновлены данные о монете")
	return result, err
}
