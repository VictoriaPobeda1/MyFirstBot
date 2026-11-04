package users

import (
	"botproject/internal/config"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInsertUserDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)
	sqlErr := errors.New("ошибка выполнения SQL")
	testsData := []struct {
		nameTest    string
		id_user     string
		time_update int
		auto_update bool
		expected    error
	}{
		{
			nameTest:    "Успешная вставка",
			id_user:     "234567090",
			time_update: 1,
			auto_update: true,
			expected:    nil,
		},
		{
			nameTest:    "Пустой id",
			id_user:     "",
			time_update: 1,
			auto_update: true,
			expected:    sqlErr,
		},
		{
			nameTest:    "Неуникальный пользователь",
			id_user:     "234567890",
			time_update: 1,
			auto_update: true,
			expected:    sqlErr,
		},
	}
	for _, testUnit := range testsData {
		t.Run(testUnit.id_user, func(t *testing.T) {
			result := newUserRepo.InsertUserDB(testUnit.id_user, testUnit.time_update, testUnit.auto_update)
			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}

func TestSelectAllUsersDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)

	testsData := []struct {
		nameTest string
		expected error
	}{
		{
			nameTest: "Успешное получение пользователей",
			expected: nil,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newUserRepo.SelectAllUsersDB()

			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}
func TestSelectAllUsersByAutoUpdate(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)

	testsData := []struct {
		nameTest string
		expected error
	}{
		{
			nameTest: "Получение пользователей с автообновлением",
			expected: nil,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newUserRepo.SelectAllUsersByAutoUpdate()

			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}
func TestSelectUserByIdDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)
	sqlErr := errors.New("ошибка получения пользователя")

	testsData := []struct {
		nameTest string
		id_user  string
		expected error
	}{
		{
			nameTest: "Существующий пользователь",
			id_user:  "490311952",
			expected: nil,
		},
		{
			nameTest: "Несуществующий пользователь",
			id_user:  "999999999",
			expected: sqlErr,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newUserRepo.SelectUserByIdDB(testUnit.id_user)

			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}
func TestUpdateTimeFromUserByIdDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)

	testsData := []struct {
		nameTest    string
		time_update int
		auto_update bool
		id_user     string
		expected    error
	}{
		{
			nameTest:    "Включение автообновления",
			time_update: 5,
			auto_update: true,
			id_user:     "234567890",
			expected:    nil,
		},
		{
			nameTest:    "Выключение автообновления",
			time_update: 5,
			auto_update: false,
			id_user:     "234567890",
			expected:    nil,
		},
		{
			nameTest:    "Несуществующий пользователь",
			time_update: 5,
			auto_update: true,
			id_user:     "999999999",
			expected:    nil,
		},
	}
	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			result := newUserRepo.UpdateTimeFromUserByIdDB(
				testUnit.time_update,
				testUnit.auto_update,
				testUnit.id_user,
			)

			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}
func TestUpdateTimeSendAtFromUserByIdDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)

	testsData := []struct {
		nameTest     string
		time_send_at time.Time
		id_user      string
		expected     error
	}{
		{
			nameTest:     "Успешное обновление времени",
			time_send_at: time.Now(),
			id_user:      "234567890",
			expected:     nil,
		},
		{
			nameTest:     "Несуществующий пользователь",
			time_send_at: time.Now(),
			id_user:      "999999999",
			expected:     nil,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			result := newUserRepo.UpdateTimeSendAtFromUserByIdDB(
				testUnit.time_send_at,
				testUnit.id_user,
			)

			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}
func TestUpdateCoinFromUserByIdDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newUserRepo := NewUserRepo(pool)

	testsData := []struct {
		nameTest string
		id_coin  int
		id_user  string
		expected error
	}{
		{
			nameTest: "Успешное обновление монеты",
			id_coin:  77,
			id_user:  "490311952",
			expected: nil,
		},
		{
			nameTest: "Несуществующий пользователь",
			id_coin:  1,
			id_user:  "999999999",
			expected: nil,
		},
	}
	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			result := newUserRepo.UpdateCoinFromUserByIdDB(
				testUnit.id_coin,
				testUnit.id_user,
			)

			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.expected,
				)
			}
		})
	}
}
