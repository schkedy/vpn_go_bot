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
	FSM     *dialog.FSMContext
	CHG     handler.CommandHandlerGroup
	CbHG    handler.CallbackHandlerGroup
	MHG     handler.MessageHandlerGroup
}

// Регистрируем диалог, который тоже участвуют в роутинге
func (r *Router) RegistrateDialog(dlg *dialog.Dialog) {
	r.Dialog = dlg
	r.CbHG.RegisterHandlers(dlg.GetHandlers())
}

// TODO: Provide storage as interface instead of dirrect redis
func (r *Router) InitializeFSMContext(ctx context.Context, userID int, storage *dialog.RedisClient) error {
	res, err := dialog.NewFSMContext(ctx, userID, r.Dialog.States, storage)
	if err != nil {
		return err
	}
	r.FSM = res
	return nil
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

// TODO: #13 сделать отдельном dialog.data storage
func InitDialogSession(ctx context.Context, FSMContext *dialog.FSMContext, userID int, storage *dialog.RedisClient) (*dialog.DialogSession, error) {
	state := FSMContext.GetState()
	return &dialog.DialogSession{
		State: state,
	}
}

func (r *Router) dispatch(update tgbotapi.Update) {
	data := &handler.HandlerData{
		Bot:    r.Bot,
		States: r.States,
	}
	ctx := context.Background()
	dialogManager := dialog.NewDialogManager(r.Dialog)

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
