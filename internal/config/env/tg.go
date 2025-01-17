package env

import (
	"crypto_scam/internal/config"
	"crypto_scam/pkg/utils/string_utils"
	"errors"
	"fmt"
	"os"
)

var _ config.TGConfig = (*tgConfig)(nil)

const (
	tgBotTokenName = "TG_BOT_TOKEN"
	tgChatIdName   = "TG_CHAT_ID"
)

type tgConfig struct {
	token  string
	chatId int64
}

// NewTgConfig создаёт новую конфигурацию взаимодействия с telegram-api в зависимости от переменных окружения
func NewTgConfig() (*tgConfig, error) {
	// получаем данные из переменных окружения
	token := os.Getenv(tgBotTokenName)
	if len(token) == 0 {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/tg.go - env variable %s not found", tgBotTokenName))
	}
	chatId := os.Getenv(tgChatIdName)
	if len(chatId) == 0 {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/tg.go - env variable %s not found", tgChatIdName))
	}
	// парсим id в int64
	id, err := string_utils.ParseID(chatId)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("internal/config/env/tg.go - env variable %s is not integer", tgChatIdName))
	}
	// возвращаем конфиг
	return &tgConfig{
		token:  token,
		chatId: id,
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
