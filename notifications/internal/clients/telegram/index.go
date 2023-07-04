package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type St struct {
	chatId int64

	bot *tgbotapi.BotAPI
}

func New(
	token string,
	chatId int64,
) (*St, error) {
	var err error

	res := &St{
		chatId: chatId,
	}

	if token == "" {
		return nil, fmt.Errorf("no token provided")
	}

	res.bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	res.bot.Debug = false

	return res, nil
}

func (o *St) Send(msg string) error {
	_, err := o.bot.Send(tgbotapi.NewMessage(o.chatId, msg))
	if err != nil {
		return fmt.Errorf("telegram send error: %w", err)
	}
	return nil
}
