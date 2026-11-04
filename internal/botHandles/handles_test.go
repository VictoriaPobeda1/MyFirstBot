package bothandles

import (
	"botproject/internal/config"
	"errors"
	"testing"

	tele "gopkg.in/telebot.v3"
)

type startTestContext struct {
	tele.Context
}

func (c *startTestContext) Sender() *tele.User {
	return &tele.User{
		ID: 234567890,
	}
}

func (c *startTestContext) Send(
	what interface{},
	opts ...interface{},
) error {
	return nil
}

func TestStart(t *testing.T) {
	ctx := &startTestContext{}
	cfg := config.Load()
	handleErr := errors.New("Ошибка подключения к API")
	testsData := []struct {
		nameTest string
		ApiBaseUrl string
		expected error
	}{
		{
			nameTest: "Успешное выполнение команды",
			ApiBaseUrl: cfg.APIBaseURL,
			expected: nil,
		},
		{
			nameTest: "Ошибка выполнения команды",
			ApiBaseUrl: "hthhhtp080",
			expected: handleErr,
		},
	}
	for _, testUnit := range testsData {
		t.Setenv("API_BASE_URL", testUnit.ApiBaseUrl)
		cfg := config.Load()
		handle := NewHandleRepo(cfg)
		result := handle.Start(ctx)
		t.Run(testUnit.nameTest, func(t *testing.T) {
			if (result != nil) != (testUnit.expected != nil) {
				t.Errorf(
					"получено %v, ожидалось %v",
					result,
					testUnit.nameTest,
				)
			}
		})
	}
}
