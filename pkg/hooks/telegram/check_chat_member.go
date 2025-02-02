package telegram

import (
	"crypto_scam/internal/config"
	"crypto_scam/internal/config/env"
	"crypto_scam/internal/logger"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var TgConfig config.TGConfig

// init инициализирует конфиг для работы с телеграм апи.
func init() {
	var err error
	// считываем конфиг
	TgConfig, err = env.NewTgConfig()
	if err != nil {
		logger.Fatal("check_chat_member.go/init - failed to create TgConfig: ", err)
	}
}

// IfChatMember проверяет, есть ли пользователь в чате приватки.
//
// Входным параметром функции является идентификатор
// пользователя, наличие в приватке которого нужно проверить.
//
// Выходными параметрами функции являются булевый флаг,
// который показывает, есть ли польщователь в чате, и ошибка,
// если она возникла, в противном случае, вместо
// неё будет возвращён nil.
func IfChatMember(id int) (bool, error) {
	// получаем токен телеграм бота
	token := TgConfig.Token()

	// создаём и инициализируем телеграм бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Fatal("check_chat_member.go/IfChatMember - failed to create bot: ", err)
		return false, err
	}

	// проверяем наличие пользователя в чате приватки
	_, err = bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             TgConfig.ChatId(),
		SuperGroupUsername: " ",
		UserID:             id,
	})

	// если возвращена ошибка.
	// пользователя нет в чате
	if err != nil {
		return false, nil
	}

	// возвращаем nil,
	// если не возникло ошибок
	return true, nil
}
