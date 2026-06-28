package dialog

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RadioItem struct {
	ID    int
	Value string
}

type RadioOnItemClick func(ctx context.Context, data *handler.HandlerData, update *tgbotapi.Update, dialogManager *DialogManager, item RadioItem)

type RadioOnStateChanged func(ctx context.Context, data *handler.HandlerData, update *tgbotapi.Update, previousItemID int, currentItemID int)

// Radio — stateful widget по аналогии с aiogram-dialog Radio:
// при нажатии на item делает его checked, снимая checked с предыдущего.
type Radio struct {
	ID string

	checkedTextTemplate   *template.Template
	uncheckedTextTemplate *template.Template

	Items string
	Width int

	onItemClick    RadioOnItemClick
	onStateChanged RadioOnStateChanged
	dialogManager  *DialogManager

	mu            sync.RWMutex
	checkedItemID int
}

func NewRadio(
	id string,
	checkedText string,
	uncheckedText string,
	items string,
	width int,
	onItemClick RadioOnItemClick,
	onStateChanged RadioOnStateChanged,
) *Radio {
	if width <= 0 {
		width = 1
	}

	checkedTmpl, _ := template.New("radio_checked").Parse(checkedText)
	uncheckedTmpl, _ := template.New("radio_unchecked").Parse(uncheckedText)

	return &Radio{
		ID:                    id,
		checkedTextTemplate:   checkedTmpl,
		uncheckedTextTemplate: uncheckedTmpl,
		Items:                 items,
		Width:                 width,
		onItemClick:           onItemClick,
		onStateChanged:        onStateChanged,
		checkedItemID:         -1,
	}
}

func (r *Radio) GetChecked() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.checkedItemID
}

func (r *Radio) IsChecked(itemID int) bool {
	return r.GetChecked() == itemID
}

func (r *Radio) SetChecked(itemID int) {
	r.mu.Lock()
	r.checkedItemID = itemID
	r.mu.Unlock()
}

func (r *Radio) SetDialogManager(dialogManager *DialogManager) {
	r.dialogManager = dialogManager
}

func (r *Radio) getButtonRows(data map[string]interface{}) []ButtonRow {
	if r == nil {
		return nil
	}

	items := r.resolveItems(data)
	if len(items) == 0 {
		return nil
	}

	rows := make([]ButtonRow, 0, (len(items)+r.Width-1)/r.Width)
	currentRow := make([]*Button, 0, r.Width)

	for _, item := range items {
		buttonText := r.renderItemText(data, item)
		button := NewButton(buttonText, r.callbackData(item.ID), nil, "")
		currentRow = append(currentRow, button)

		if len(currentRow) == r.Width {
			rows = append(rows, ButtonRow{Buttons: currentRow})
			currentRow = make([]*Button, 0, r.Width)
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, ButtonRow{Buttons: currentRow})
	}

	return rows
}

func (r *Radio) getHandlers() map[string]handler.HandlerFunc {
	handlers := make(map[string]handler.HandlerFunc)
	if r == nil {
		return handlers
	}

	callbackPrefix := r.callbackPrefix()
	handlers[callbackPrefix] = func(ctx context.Context, data *handler.HandlerData, update *tgbotapi.Update) {
		itemID := r.parseItemID(update)
		if itemID < 0 {
			return
		}

		items := r.resolveItems(r.resolveHandlerData(data))
		selectedItem, ok := r.findItemByID(items, itemID)
		if !ok {
			selectedItem = RadioItem{ID: itemID}
		}

		previousItemID := r.GetChecked()
		r.SetChecked(selectedItem.ID)

		if r.onItemClick != nil {
			r.onItemClick(ctx, data, update, r.dialogManager, selectedItem)
		}

		if previousItemID != selectedItem.ID && r.onStateChanged != nil {
			r.onStateChanged(ctx, data, update, previousItemID, selectedItem.ID)
		}
	}

	return handlers
}

func (r *Radio) callbackData(itemID int) string {
	return fmt.Sprintf("radio:%s:%d", r.ID, itemID)
}

func (r *Radio) callbackPrefix() string {
	return fmt.Sprintf("radio:%s:", r.ID)
}

func (r *Radio) parseItemID(update *tgbotapi.Update) int {
	if update == nil || update.CallbackQuery == nil {
		return -1
	}

	data := update.CallbackQuery.Data
	prefix := r.callbackPrefix()
	if len(data) <= len(prefix) || data[:len(prefix)] != prefix {
		return -1
	}

	itemID, err := strconv.Atoi(data[len(prefix):])
	if err != nil {
		return -1
	}

	return itemID
}

func (r *Radio) findItemByID(items []RadioItem, itemID int) (RadioItem, bool) {
	for _, item := range items {
		if item.ID == itemID {
			return item, true
		}
	}
	return RadioItem{}, false
}

func (r *Radio) resolveHandlerData(handlerData *handler.HandlerData) map[string]interface{} {
	if handlerData == nil || handlerData.GetterData == nil {
		return map[string]interface{}{}
	}

	return handlerData.GetterData
}

func (r *Radio) renderItemText(data map[string]interface{}, item RadioItem) string {
	tmpl := r.uncheckedTextTemplate
	if r.IsChecked(item.ID) {
		tmpl = r.checkedTextTemplate
	}

	if tmpl == nil {
		return ""
	}

	templateData := make(map[string]interface{}, len(data)+3)
	for key, value := range data {
		templateData[key] = value
	}
	templateData["item"] = item.Value
	templateData["item_id"] = item.ID
	templateData["checked"] = r.IsChecked(item.ID)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return ""
	}

	return buf.String()
}

func (r *Radio) resolveItems(data map[string]interface{}) []RadioItem {
	if r == nil || r.Items == "" || data == nil {
		return nil
	}

	rawItems, exists := data[r.Items]
	if !exists || rawItems == nil {
		return nil
	}

	items, ok := rawItems.([]RadioItem)
	if !ok {
		return nil
	}

	return items
}
