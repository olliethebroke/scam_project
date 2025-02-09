package telegram

import (
	"crypto_scam/internal/config"
	"crypto_scam/internal/logger"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var tgConfig config.TGConfig
var bot *tgbotapi.BotAPI

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
	// проверяем наличие пользователя в чате приватки
	_, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             tgConfig.ChatId(),
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

// InitTGConfig инициализирует переменные tgConfig и bot.
func InitTGConfig(cfg config.TGConfig) {
	tgConfig = cfg

	// получаем токен телеграм бота
	token := tgConfig.Token()

	// создаём и инициализируем телеграм бота
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Fatal("check_chat_member.go/IfChatMember - failed to create bot: ", err)
	}
}

// TGConfig возвращает конфиг взаимодействия с тг.
func TGConfig() config.TGConfig {
	return tgConfig
}
