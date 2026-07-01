package dialog

import "context"

// Widget описывает любой UI-элемент, который может отрисовать кнопки и вернуть связанные обработчики.
type Widget interface {
	getButtonRows(ctx context.Context, fsm FSMContext, data map[string]interface{}) []ButtonRow
	getHandlers() map[string]HandlerFunc
}

type dialogManagerAwareWidget interface {
	SetDialogManager(dialogManager *DialogManager)
}

// MediaProvider — дополнительный интерфейс виджета, который умеет возвращать текущее медиа для рендера окна.
// Window проверяет виджеты на его реализацию и, если находит, использует медиа от виджета.
type MediaProvider interface {
	GetCurrentMedia(ctx context.Context, fsm FSMContext, data map[string]interface{}) *Media
}
