package dialog

import (
	"context"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type HandlerFunc func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update)

type InnerHandler func(ctx context.Context, getterData map[string]interface{}, session *DialogSession, update *tgbotapi.Update)
