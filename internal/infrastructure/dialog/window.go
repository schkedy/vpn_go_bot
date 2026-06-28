package dialog

import (
	"bytes"
	"text/template"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// Как можно реализовать widget типа Radio как в airogram_dialog, чтобы он был совмещен с моим проектом
type Window struct {
	State        State              // состояние, при котором открывается окно
	textTemplate *template.Template // скомпилированный шаблон для динамического текста
	Getter       Getter             // getter всегда получает DialogManager и deps
	widgets      []Widget           // виджеты, которые будут отображаться в окне
	media        *Media             // медиа-файл (фото, видео, аудио и т.д.)
}

// NewWindow to make Window
func NewWindow(state State, text string, getter Getter, widgets ...Widget) *Window {
	// Компилируем текст как шаблон
	tmpl, _ := template.New("text").Parse(text)

	w := &Window{
		State:        state,
		textTemplate: tmpl,
		Getter:       getter,
		widgets:      widgets,
	}
	return w
}

// WindowConfig is  result of Window Rendering
type WindowConfig struct {
	Text     string
	Keyboard tgbotapi.InlineKeyboardMarkup
	Media    *Media
}

// RenderAll as collection point, render all window and pass dialogManager to widgets
func (w *Window) RenderAll(dialogManager *DialogManager) WindowConfig {
	data := w.getGetterData(dialogManager)
	text := w.renderText(data)
	keyboard := w.renderWidgets(data)
	media := w.renderMedia(data)

	return WindowConfig{
		Text:     text,
		Keyboard: keyboard,
		Media:    media,
	}
}

// Render dynamic text
// ? как это работает
func (w *Window) renderText(data map[string]interface{}) string {
	if w.textTemplate == nil {
		return ""
	}

	var buf bytes.Buffer
	if err := w.textTemplate.Execute(&buf, data); err != nil {
		return ""
	}

	return buf.String()
}

// renderMedia collect media from widgets  first and Window.media second
func (w *Window) renderMedia(data map[string]interface{}) *Media {
	for _, widget := range w.widgets {
		if provider, ok := widget.(MediaProvider); ok {
			if m := provider.GetCurrentMedia(data); m != nil {
				return m
			}
		}
	}

	if w.shouldRenderMedia(data) {
		return w.media
	}

	return nil
}

// renderWidgets collect all widgets into one InlineKeyboardMarkup
func (w *Window) renderWidgets(data map[string]interface{}) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()

	for _, widget := range w.widgets {
		for _, row := range widget.getButtonRows(data) {
			keyboardRow := make([]tgbotapi.InlineKeyboardButton, 0, len(row.Buttons))
			for _, button := range row.Buttons {
				renderedButton := button.ToInlineKeyboardButton(data)
				if renderedButton != nil {
					keyboardRow = append(keyboardRow, *renderedButton)
				}
			}

			if len(keyboardRow) > 0 {
				markup.InlineKeyboard = append(markup.InlineKeyboard, keyboardRow)
			}
		}
	}

	return markup
}

func (w *Window) shouldRenderMedia(data map[string]interface{}) bool {
	if w.media == nil || w.media.HasMediaKey == "" {
		return false
	}

	if val, exists := data[w.media.HasMediaKey]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}

		if strVal, ok := val.(string); ok {
			return strVal != ""
		}
	}

	return false
}

type WindowGetter func(dialogManager *DialogManager, deps map[string]interface{}) map[string]interface{}

// GetHandler return map of widget handlers
func (w *Window) GetHandlers() map[string]HandlerFunc {
	handlers := make(map[string]HandlerFunc)
	for _, widget := range w.widgets {
		for key, handlerFunc := range widget.getHandlers() {
			handlers[key] = handlerFunc
		}
	}

	return handlers
}

// BindDialogManager provide DialogManager to every widget
func (w *Window) BindDialogManager(dialogManager *DialogManager) {
	for _, widget := range w.widgets {
		if bindableWidget, ok := widget.(dialogManagerAwareWidget); ok {
			bindableWidget.SetDialogManager(dialogManager)
		}
	}
}

// Collect getterData
func (w *Window) getGetterData(dialogManager *DialogManager) map[string]interface{} {
	var data map[string]interface{}

	if w.Getter != nil {
		data = w.Getter(dialogManager.Session, dialogManager.deps)
	}

	if data == nil {
		return map[string]interface{}{}
	}

	return data
}
