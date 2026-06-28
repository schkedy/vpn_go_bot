package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

const (
	token    = "8236311725:AAGUOA_IF1fiE1HlgumczRanuUYCXuBlEOE"
	photoURL = "https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_1280.jpg"
	apiBase  = "https://api.telegram.org/bot" + token
)

var (
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Button 1", "button1"),
			tgbotapi.NewInlineKeyboardButtonData("Button 2", "button2"),
		),
	)
)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(photoURL))
		// photoMsg.Caption = "Here is a photo!"
		// photoMsg.ReplyMarkup = keyboard
		Msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here is a message with buttons!")
		sentPhoto, err := bot.Send(Msg)
		if err != nil {
			log.Printf("send photo message error: %v", err)
			continue
		}
		time.Sleep(2 * time.Second)

		media := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_1280.jpg"))
		media.Caption = "Updated caption with buttons!"
		_, err = bot.EditMessageMediaAndMarkup(
			update.Message.Chat.ID,
			sentPhoto.MessageID,
			&media,
			keyboard,
		)
		if err != nil {
			log.Printf("edit message media error: %v", err)
			continue
		}
		time.Sleep(2 * time.Second)
		_, err = bot.EditMessageCaption(
			update.Message.Chat.ID,
			sentPhoto.MessageID,
			"New text with buttons!",
		)
		if err != nil {
			log.Printf("edit message text error: %v", err)
			continue
		}
	}
}
