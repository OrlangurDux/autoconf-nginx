package bot

import (
	"autoconf/config"
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
	//ChatID telegram chat id
	ChatID int64
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
		reply := "Command don't find."
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "/join KEY - for joined monitoring service autoconf"
		case "join":
			key := strings.Split(update.Message.Text, " ")
			if key[1] == JoinKey {
				reply = "Joined monitoring"
				config.SysConfig.Telegram.ChatID = update.Message.Chat.ID
				config.WriteSysConfig()
			} else {
				reply = "Error joined monitoring"
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		log.Println(update.Message.Chat.ID)
		Bot.Send(msg)
	}
}

//SendBotMessage send in telegram bot message
func SendBotMessage(msg string) {
	message := tgbotapi.NewMessage(ChatID, msg)
	Bot.Send(message)
}
