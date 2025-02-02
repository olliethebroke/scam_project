package env

import (
	"crypto_scam/internal/config"
	"crypto_scam/internal/utils/string_utils"
	"errors"
	"fmt"
	"os"
	"time"
)

var _ config.TGConfig = (*tgConfig)(nil)

const (
	tgBotTokenName  = "TG_BOT_TOKEN"
	tgChatIdName    = "TG_CHAT_ID"
	initDataExpName = "INIT_DATA_EXP_TIME=3"
)

// tgConfig стуктура-конфигурация взаимодействия с telegram-api
// token - токен бота
// chatId - id чата приватки
// initDataExp - время действия данных
type tgConfig struct {
	token       string
	chatId      int64
	initDataExp time.Duration
}

// NewTgConfig создаёт новую конфигурацию взаимодействия с telegram-api в зависимости от переменных окружения
func NewTgConfig() (*tgConfig, error) {
	// получаем данные из переменных окружения
	token := os.Getenv(tgBotTokenName)
	if len(token) == 0 {
		return nil,
			errors.New(fmt.Sprintf("tg.go/NewTgConfig - env variable %s not found", tgBotTokenName))
	}
	chatId := os.Getenv(tgChatIdName)
	if len(chatId) == 0 {
		return nil,
			errors.New(fmt.Sprintf("tg.go/NewTgConfig - env variable %s not found", tgChatIdName))
	}
	initDataExpStr := os.Getenv(initDataExpName)
	if len(initDataExpStr) == 0 {
		return nil,
			errors.New(fmt.Sprintf("tg.go/NewTgConfig - env variable %s not found", initDataExpName))
	}
	// парсим id в int64
	id, err := string_utils.StringToInt64(chatId)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("tg.go/NewTgConfig - env variable %s is not integer", tgChatIdName))
	}
	// парсим initDataExpStr в int64
	initDataExp, err := string_utils.StringToInt64(initDataExpStr)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("tg.go/NewTgConfig - env variable %s is not integer", initDataExpName))
	}
	// возвращаем конфиг
	return &tgConfig{
		token:       token,
		chatId:      id,
		initDataExp: time.Duration(initDataExp) * time.Minute,
	}, nil
}

// Token возвращает токен телеграм бота
func (cfg *tgConfig) Token() string {
	return cfg.token
}

// ChatId возвращает id чата, с которым будет работать телеграм бот
func (cfg *tgConfig) ChatId() int64 {
	return cfg.chatId
}

func (cfg *tgConfig) InitDataExpiration() time.Duration {
	return cfg.initDataExp
}
