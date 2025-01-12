package telegram

import (
	"crypto_scam/internal/config/env"
	"crypto_scam/internal/logger/logrus_logger"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// IfChatMember проверяет, есть ли пользователь в чате приватки
func IfChatMember(id int) error {
	tgConfig, err := env.NewTgConfig()
	if err != nil {
		logrus_logger.Log.Fatal("pkg/hooks/telegram/check_chat_member.go - failed to create tgConfig")
	}
	token := tgConfig.Token()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus_logger.Log.Fatal("pkg/hooks/telegram/check_chat_member.go - failed to create bot")
	}
	_, err = bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             tgConfig.ChatId(),
		SuperGroupUsername: " ",
		UserID:             id,
	})
	return err
}
