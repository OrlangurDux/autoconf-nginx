package bot

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	//BotToken telegram bot token
	BotToken string
	//JoinKey key for connect notification
	JoinKey string
	//Bot for boot api
	Bot *tgbotapi.BotAPI
)

func InitBot() {
	var err error
	for {
		Bot, err = tgbotapi.NewBotAPI(BotToken)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Millisecond * time.Duration(10000))
		} else {
			break
		}
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := Bot.GetUpdatesChan(u)

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
				reply = "Joined monitoring"
			} else {
				reply = "Error joined monitoring"
			}
		}
		//log.Println(update.Message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		log.Println(update.Message.Chat.ID)

		Bot.Send(msg)
	}
}

//SendBotMessage send in telegram bot message
func SendBotMessage(msg string) {

	message := tgbotapi.NewMessage(355199786, msg)

	Bot.Send(message)
}
