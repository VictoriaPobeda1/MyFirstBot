package users

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserFunc interface {
	InsertUserDB(id_user string, time_update int, auto_update bool) error
	SelectAllUsersDB() ([]User, error)
	SelectAllUsersByAutoUpdate() ([]UserCoin, error)
	SelectUserByIdDB(id_user string) (User, error)
	UpdateTimeFromUserByIdDB(time_update int, auto_update bool, id_user string) error
	UpdateTimeSendAtFromUserByIdDB(time_send_at time.Time, id_user string) error
	UpdateCoinFromUserByIdDB(id_coin int, id_user string) error
}

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (ur *UserRepo) InsertUserDB(id_user string, time_update int, auto_update bool) error {
	sql := `insert into "CriptoProject".users(
	id_user, time_update, auto_update)
	values ($1,$2,$3)`
	_, err := ur.pool.Exec(context.Background(), sql, id_user, time_update, auto_update)
	if err != nil {
		slog.Error("Ошибка добавления пользователя в БД:", err)
	}
	slog.Info("Пользователь добавлен")
	return err
}
func (ur *UserRepo) SelectAllUsersDB() ([]User, error) {
	sql := `select id_user, time_update, auto_update, id_coin
	from "CriptoProject".users`
	result, err := ur.pool.Query(context.Background(), sql)
	if err != nil {
		slog.Error("Ошибка вывода пользователей из БД:", err)
	}
	defer result.Close()
	users := []User{}
	for result.Next() {
		u := User{}
		err := result.Scan(&u.Id_user, &u.Time_update, &u.Auto_update, &u.Id_coin)
		if err != nil {
			slog.Error("Ошибка считывания данных о пользователе:", err)
			continue
		}
		users = append(users, u)
	}
	slog.Info("Данные о пользователях выведены из БД")
	return users, err
}
func (ur *UserRepo) SelectAllUsersByAutoUpdate() ([]UserCoin, error) {
	sql := `select "CriptoProject".coin.*,"CriptoProject".users.auto_update,
	"CriptoProject".users.time_update,"CriptoProject".users.id_user,"CriptoProject".users.time_send_at 
	from "CriptoProject".users inner join "CriptoProject".coin
	on "CriptoProject".users.id_coin="CriptoProject".coin.id
	where "CriptoProject".users.auto_update=true`
	result, err := ur.pool.Query(context.Background(), sql)
	if err != nil {
		slog.Error("Ошибка вывода пользователей и их монет из БД по условию autoupdate:", err)
	}
	defer result.Close()
	userCoin := []UserCoin{}
	for result.Next() {
		u := UserCoin{}
		err := result.Scan(&u.Id_coin, &u.Name, &u.Price, &u.Percent_change_1h, &u.Price_min, &u.Price_max, &u.Auto_update, &u.Time_update, &u.Id_user, &u.Time_send_at)
		if err != nil {
			slog.Error("Ошибка считывания данных пользователей и их монет:", err)
			continue
		}
		userCoin = append(userCoin, u)
	}
	slog.Info("Пользователи выведены")
	return userCoin, err
}
func (ur *UserRepo) SelectUserByIdDB(id_user string) (User, error) {
	sql := `select id_user, time_update, auto_update, id_coin
	from "CriptoProject".users where id_user=$1`
	user := User{}
	err := ur.pool.QueryRow(context.Background(), sql, id_user).Scan(&user.Id_user, &user.Time_update, &user.Auto_update, &user.Id_coin)
	if err != nil {
		slog.Error("Ошибка вывода пользователя из БД по id:", err)
	}
	slog.Info("Пользователь выведен")
	return user, err
}
func (ur *UserRepo) UpdateTimeFromUserByIdDB(time_update int, auto_update bool, id_user string) error {
	sql := `update "CriptoProject".users
	set time_update=$1, auto_update=$2,time_send_at=$3
	where id_user=$4`
	_, err := ur.pool.Exec(context.Background(), sql, time_update, auto_update, time.Now(), id_user)
	if err != nil {
		slog.Error("Ошибка обновления данных о пользователе в БД:", err)
	}
	slog.Info("Обновлены настройки автоотправки пользователя")
	return err
}
func (ur *UserRepo) UpdateTimeSendAtFromUserByIdDB(time_send_at time.Time, id_user string) error {
	sql := `update "CriptoProject".users
	set time_send_at=$1 where id_user=$2`
	_, err := ur.pool.Exec(context.Background(), sql, time_send_at, id_user)
	if err != nil {
		slog.Error("Ошибка обновления времени последней отправки у пользователя в БД:", err)
	}
	slog.Info("Обновлено время последней отправки у пользователя")
	return err
}
func (ur *UserRepo) UpdateCoinFromUserByIdDB(id_coin int, id_user string) error {
	sql := `update "CriptoProject".users
	set id_coin=$1
	where id_user=$2`
	_, err := ur.pool.Exec(context.Background(), sql, id_coin, id_user)
	if err != nil {
		slog.Error("Ошибка обновления данных о пользователе в БД:", err)
	}
	slog.Info("Монета у пользователя обновлена")
	return err
}
