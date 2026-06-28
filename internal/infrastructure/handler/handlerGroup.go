package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State string

type HandlerGroup interface {
	Validate(s string) (command string)
	GetHandler(command string) (HandlerFunc, bool)
}

type HandlerData struct {
	Bot    *tgbotapi.BotAPI
	States *[]State
	GetterData map[string]interface{}
	// add services, repositories, etc
}

type HandlerFunc func(ctx context.Context, data *HandlerData, update *tgbotapi.Update)
