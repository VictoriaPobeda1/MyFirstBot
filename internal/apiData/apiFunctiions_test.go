package apidata

import (
	config "botproject/internal/config"
	"errors"
	"testing"
)

func TestGetAllCoinFromApi(t *testing.T) {
	api := &ApiRepo{}
	cfg := config.Load()
	apiErr := errors.New("Ошибка API")
	testsData := []struct {
		nameTest string
		tokenApi string
		expected error
	}{
		{
			nameTest: "Успешный вывод монеты",
			tokenApi: cfg.API_KEY,
			expected: nil,
		},
		{
			nameTest: "Ошибка API",
			tokenApi: "fghjhgfdfgh",
			expected: apiErr,
		},
	}
	for _, testUnit := range testsData {
		t.Setenv("API_KEY", testUnit.tokenApi)
		result := api.GetAllCoinFromApi()
		t.Run(testUnit.nameTest, func(t *testing.T) {
			if (len(result.Data) != 0) != (testUnit.expected == nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.nameTest,
				)
			}
		})
	}
}
