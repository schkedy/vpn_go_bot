package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandlerGroup struct {
}

func (ch *MessageHandlerGroup) Validate(s string) (command string) {
	return ""
}

func (ch *MessageHandlerGroup) GetHandler(command string) (HandlerFunc, bool) {
	return messageHandler, true
}

func messageHandler(ctx context.Context, data *HandlerData, update *tgbotapi.Update) {
	// handle non-command messages
}
