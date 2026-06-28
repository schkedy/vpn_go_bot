package dialog

// Widget описывает любой UI-элемент, который может отрисовать кнопки и вернуть связанные обработчики.
type Widget interface {
	getButtonRows(data map[string]interface{}) []ButtonRow
	getHandlers() map[string]HandlerFunc
}

type dialogManagerAwareWidget interface {
	SetDialogManager(dialogManager *DialogManager)
}

// MediaProvider — дополнительный интерфейс виджета, который умеет возвращать текущее медиа для рендера окна.
// Window проверяет виджеты на его реализацию и, если находит, использует медиа от виджета.
type MediaProvider interface {
	GetCurrentMedia(data map[string]interface{}) *Media
}
