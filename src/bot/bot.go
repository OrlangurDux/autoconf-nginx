package bot

import (
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	//BotToken telegram bot token
	BotToken string
)

//SendBotMessage send in telegram bot message
func SendBotMessage(msg string) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Print(err)
	}

	message := tgbotapi.NewMessage(355199786, msg)

	bot.Send(message)
}
