package bot

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	//BotToken telegram bot token
	BotToken string
	//JoinKey key for connect notification
	JoinKey string
	bot     *tgbotapi.BotAPI
)

func InitBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// u - struct with config for get update
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// use config u create channel for push message
	updates, _ := bot.GetUpdatesChan(u)
	// in channel updates send truct type Update
	// reading and processing request
	for update := range updates {
		log.Println(update.Message.Command())
		reply := "I don't now"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Hi. I'am telegram bot"
		case "join":
			key := strings.Split(update.Message.Text, " ")
			if key[1] == JoinKey {
				reply = "Joined"
			} else {
				reply = "Error joined channel"
			}
		}
		//log.Println(update.Message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		log.Println(update.Message.Chat.ID)

		bot.Send(msg)
	}
}

//SendBotMessage send in telegram bot message
func SendBotMessage(msg string) {

	message := tgbotapi.NewMessage(355199786, msg)

	bot.Send(message)
}
