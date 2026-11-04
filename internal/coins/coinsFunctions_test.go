package coins

import (
	"botproject/internal/config"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInsertCoinDB(t *testing.T) {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		t.Fatalf("создание пула: %v", err)
	}
	defer pool.Close()
	newCoinRepo := NewCoinRepo(pool)
	sqlErr := errors.New("ошибка выполнения SQL")
	testsData := []struct {
		nameTest          string
		name              string
		price             float64
		percent_change_1h float64
		price_min         float64
		price_max         float64
		expected          error
	}{
		{
			nameTest:          "Успешная вставка",
			name:              "Bitcoin",
			price:             100,
			percent_change_1h: 1,
			price_min:         90,
			price_max:         110,
			expected:          nil,
		},
		{
			nameTest:          "Отрицательная цена",
			name:              "Bitcoin",
			price:             -100,
			percent_change_1h: 1,
			price_min:         90,
			price_max:         110,
			expected:          sqlErr,
		},
	}
	for _, testUnit := range testsData {
		t.Run(testUnit.name, func(t *testing.T) {
			result := newCoinRepo.InsertCoinDB(testUnit.name, testUnit.price, testUnit.percent_change_1h, testUnit.price_min, testUnit.price_max)
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
func TestSelectAllCoinsDB(t *testing.T) {
	cfg := config.Load()
	pool, _ := pgxpool.New(context.Background(), cfg.DSN)
	defer pool.Close()
	newCoinRepo := NewCoinRepo(pool)

	testsData := []struct {
		nameTest string
		expected error
	}{
		{
			nameTest: "Успешное получение монет",
			expected: nil,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newCoinRepo.SelectAllCoinsDB()

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
func TestSelectCoinByIdDB(t *testing.T) {
	cfg := config.Load()
	pool, _ := pgxpool.New(context.Background(), cfg.DSN)
	defer pool.Close()
	newCoinRepo := NewCoinRepo(pool)
	sqlErr := errors.New("ошибка получения монеты")

	testsData := []struct {
		nameTest string
		id       int
		expected error
	}{
		{
			nameTest: "Существующая монета",
			id:       7,
			expected: nil,
		},
		{
			nameTest: "Несуществующая монета",
			id:       -1,
			expected: sqlErr,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newCoinRepo.SelectCoinByIdDB(testUnit.id)

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
func TestSelectCoinByNameDB(t *testing.T) {
	cfg := config.Load()
	pool, _ := pgxpool.New(context.Background(), cfg.DSN)
	defer pool.Close()
	newCoinRepo := NewCoinRepo(pool)
	sqlErr := errors.New("ошибка получения монеты")

	testsData := []struct {
		nameTest string
		name     string
		expected error
	}{
		{
			nameTest: "Существующая монета",
			name:     "Bitcoin",
			expected: nil,
		},
		{
			nameTest: "Несуществующая монета",
			name:     "UnknownCoinForTest",
			expected: sqlErr,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newCoinRepo.SelectCoinByNameDB(testUnit.name)

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
func TestUpdateCoinByIdDB(t *testing.T) {
	cfg := config.Load()
	pool, _ := pgxpool.New(context.Background(), cfg.DSN)
	defer pool.Close()
	newCoinRepo := NewCoinRepo(pool)
	sqlErr := errors.New("ошибка обновления монеты")

	testsData := []struct {
		nameTest          string
		price             float64
		percent_change_1h float64
		id                int
		expected          error
	}{
		{
			nameTest:          "Успешное обновление",
			price:             100,
			percent_change_1h: 1,
			id:                7,
			expected:          nil,
		},
		{
			nameTest:          "Отрицательная цена",
			price:             -100,
			percent_change_1h: 1,
			id:                7,
			expected:          sqlErr,
		},
	}

	for _, testUnit := range testsData {
		t.Run(testUnit.nameTest, func(t *testing.T) {
			_, result := newCoinRepo.UpdateCoinByIdDB(
				testUnit.price,
				testUnit.percent_change_1h,
				testUnit.id,
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
