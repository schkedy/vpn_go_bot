package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandlerGroup struct {
	handlers map[string]HandlerFunc
}

func NewCommandHandlers() *CommandHandlerGroup {
	handlers := make(map[string]HandlerFunc)
	handlers["start"] = StartCommandHandler

	return &CommandHandlerGroup{
		handlers: handlers,
	}
}

func (ch *CommandHandlerGroup) Validate(s string) (command string) {
	if len(s) > 0 && s[0] == '/' {
		for i, c := range s {
			if c == ' ' {
				return s[1:i]
			}
		}
		return s[1:]
	}
	return ""
}

func (chg *CommandHandlerGroup) GetHandler(command string) (HandlerFunc, bool) {
	handler, exists := chg.handlers[command]
	return handler, exists
}
func StartCommandHandler(ctx context.Context, data *HandlerData, update *tgbotapi.Update) {
	// handle command
}
