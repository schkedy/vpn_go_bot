package router

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Router struct {
	States    *[]State
	Updates   tgbotapi.UpdatesChannel
	Bot       *tgbotapi.BotAPI
	matchFunc func(update tgbotapi.Update) // TODO: #34 нужно добавить остальную нагрузку кроме updates, например, бота и состояние диалогов
}

func (r *Router) Route(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		match(update)
	}
}
