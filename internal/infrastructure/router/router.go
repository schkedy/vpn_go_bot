package router

import (
	"context"
	"vpn_go_bot/internal/infrastructure/dialog"
	"vpn_go_bot/internal/infrastructure/handler"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Router struct {
	States  *[]handler.State
	Updates tgbotapi.UpdatesChannel
	Bot     *tgbotapi.BotAPI
	Dialog  *dialog.Dialog
	CHG     handler.CommandHandlerGroup
	CbHG    handler.CallbackHandlerGroup
	MHG     handler.MessageHandlerGroup
}

// Регистрируем диалог, который тоже участвуют в роутинге
func (r *Router) RegistrateDialog(dlg *dialog.Dialog) {
	r.Dialog = dlg
	r.CbHG.RegisterHandlers(dlg.GetHandlers())
}

func (r *Router) Route(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		u := update
		go r.dispatch(u)

	}
}

// Прокидывать dialogManager на обработку каждого update
// dialogManager должен собираться из редис если есть данные, если нет новым пустым
//

// dispatch выбирает куда пойдет update
// callbackData  определяет в какой cal// getterFuncData := dm.dialog.GetWindow(dm.state).Getter(dm) - получаем данные для отображения окна
// message := dm.dialog.GetWindow(dm.state).getRenderWindow(getterFuncData) - получаем окно
// sendMessage(message) - отправляем сообщениеlbackHandler пойдет update
//

func (r *Router) dispatch(update tgbotapi.Update) {
	data := &handler.HandlerData{
		Bot:    r.Bot,
		States: r.States,
	}
	var snd dialog.Sender = &r.Bot.
	ctx := context.Background()
	switch whatType(update) {
	case "message":
		if update.Message != nil && update.Message.IsCommand() {
			cmd := update.Message.Command()
			if h, ok := r.CHG.GetHandler(cmd); ok {
				h(ctx, data, &update)
				return
			}
		}
		h, _ := r.MHG.GetHandler("")
		h(ctx, data, &update)
		return

	case "callback_query":
		if update.CallbackQuery == nil {
			return
		}
		cbData := update.CallbackQuery.Data
		if h, ok := r.CbHG.GetHandler(cbData); ok {
			h(ctx, data, &update)
			return
		}
	}
}

func whatType(update tgbotapi.Update) string {
	if update.Message != nil {
		return "message"
	}
	if update.CallbackQuery != nil {
		return "callback_query"
	}
	return "other"
}
