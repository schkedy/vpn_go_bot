package dialog

import (
	"bytes"
	"text/template"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// Button описывает кнопку с тем же принципом рендеринга текста, что и у Window:
// текст компилируется как шаблон и рендерится на данных из getter.
type Button struct {
	textTemplate *template.Template
	CallbackData string
	Handler      HandlerFunc
	ShowWhenKey  string
}

// ButtonRow — внутренняя строка кнопок для inline keyboard, каждая строка отвечает за кнопки на одной шеренге
type ButtonRow struct {
	Buttons []*Button
}

func NewButton(text string, callbackData string, handlerFunc HandlerFunc, showWhenKey string) *Button {
	tmpl, _ := template.New("button_text").Parse(text)

	return &Button{
		textTemplate: tmpl,
		CallbackData: callbackData,
		Handler:      handlerFunc,
		ShowWhenKey:  showWhenKey,
	}
}

// GetRenderedText make button text from getter data
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

// ShouldShow search ShownKey in getter data, and return bool to show button 
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
// ToInlineKeyboardButton get telegram api button type to send
func (b *Button) ToInlineKeyboardButton(data map[string]interface{}) *tgbotapi.InlineKeyboardButton {
	if !b.ShouldShow(data) {
		return nil
	}

	button := tgbotapi.NewInlineKeyboardButtonData(b.GetRenderedText(data), b.CallbackData)
	return &button
}
// 
func (b *Button) getButtonRows(data map[string]interface{}) []ButtonRow {
	if !b.ShouldShow(data) {
		return nil
	}

	return []ButtonRow{{Buttons: []*Button{b}}}
}

// getHandlers return binding HandlerFunc to concrete button
func (b *Button) getHandlers() map[string]HandlerFunc {
	if b == nil || b.CallbackData == "" || b.Handler == nil {
		return map[string]HandlerFunc{}
	}

	return map[string]HandlerFunc{
		b.CallbackData: b.Handler,
	}
}
