package dialog

import (
	"bytes"
	"text/template"
	"vpn_go_bot/internal/infrastructure/handler"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// Button описывает кнопку с тем же принципом рендеринга текста, что и у Window:
// текст компилируется как шаблон и рендерится на данных из getter.
type Button struct {
	textTemplate *template.Template
	CallbackData string
	Handler      handler.HandlerFunc
	ShowWhenKey  string
}

func NewButton(text string, callbackData string, handlerFunc handler.HandlerFunc, showWhenKey string) *Button {
	tmpl, _ := template.New("button_text").Parse(text)

	return &Button{
		textTemplate: tmpl,
		CallbackData: callbackData,
		Handler:      handlerFunc,
		ShowWhenKey:  showWhenKey,
	}
}

func (b *Button) GetRenderedText(data map[string]interface{}) string {
	if b == nil || b.textTemplate == nil {
		return ""
	}

	var buf bytes.Buffer
	if err := b.textTemplate.Execute(&buf, data); err != nil {
		return ""
	}

	return buf.String()
}

func (b *Button) ShouldShow(data map[string]interface{}) bool {
	if b == nil {
		return false
	}

	if b.ShowWhenKey == "" {
		return true
	}

	val, exists := data[b.ShowWhenKey]
	if !exists {
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v != ""
	default:
		return v != nil
	}
}

func (b *Button) ToInlineKeyboardButton(data map[string]interface{}) *tgbotapi.InlineKeyboardButton {
	if !b.ShouldShow(data) {
		return nil
	}

	button := tgbotapi.NewInlineKeyboardButtonData(b.GetRenderedText(data), b.CallbackData)
	return &button
}

func (b *Button) getButtonRows(data map[string]interface{}) []ButtonRow {
	if !b.ShouldShow(data) {
		return nil
	}

	return []ButtonRow{{Buttons: []*Button{b}}}
}

func (b *Button) getHandlers() map[string]handler.HandlerFunc {
	if b == nil || b.CallbackData == "" || b.Handler == nil {
		return map[string]handler.HandlerFunc{}
	}

	return map[string]handler.HandlerFunc{
		b.CallbackData: b.Handler,
	}
}
