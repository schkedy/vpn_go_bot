package dialog

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type RadioItem struct {
	ID    int
	Value string
}

type RadioOnItemClick func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update, item RadioItem)

type RadioOnStateChanged func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update, previousItemID int, currentItemID int)

// Radio — виджет по аналогии с aiogram-dialog Radio:
// checked-состояние хранится в FSMContext, а не в самом виджете.
type Radio struct {
	ID string

	checkedTextTemplate   *template.Template
	uncheckedTextTemplate *template.Template

	Items string
	Width int

	onItemClick    RadioOnItemClick
	onStateChanged RadioOnStateChanged
}

// NewRadio создаёт новый виджет Radio и компилирует шаблоны текста для checked/unchecked состояний.
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
	}
}

// GetChecked возвращает ID выбранного элемента из FSM, либо -1 если выбор не установлен.
func (r *Radio) GetChecked(ctx context.Context, fsm FSMContext) int {
	_ = ctx
	return r.readIntFromFSM(fsm, r.checkedItemKey(), -1)
}

// IsChecked проверяет, выбран ли элемент с переданным ID.
func (r *Radio) IsChecked(ctx context.Context, fsm FSMContext, itemID int) bool {
	return r.GetChecked(ctx, fsm) == itemID
}

// SetChecked сохраняет выбранный ID элемента в FSM.
func (r *Radio) SetChecked(ctx context.Context, dialogManager *DialogManager, itemID int) {
	r.writeIntToFSM(ctx, dialogManager, r.checkedItemKey(), itemID)
}

// getButtonRows строит строки кнопок Radio на основе элементов из getter data.
func (r *Radio) getButtonRows(ctx context.Context, fsm FSMContext, data map[string]interface{}) []ButtonRow {
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
		buttonText := r.renderItemText(ctx, fsm, data, item)
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

// getHandlers регистрирует callback-обработчик выбора элемента Radio.
func (r *Radio) getHandlers() map[string]HandlerFunc {
	handlers := make(map[string]HandlerFunc)
	if r == nil {
		return handlers
	}

	callbackPrefix := r.callbackPrefix()
	handlers[callbackPrefix] = func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update) {
		itemID := r.parseItemID(update)
		if itemID < 0 {
			return
		}

		items := r.resolveItems(r.resolveHandlerData(ctx, dialogManager))
		selectedItem, ok := r.findItemByID(items, itemID)
		if !ok {
			selectedItem = RadioItem{ID: itemID}
		}

		previousItemID := r.readIntFromDialogManagerFSM(ctx, dialogManager, r.checkedItemKey(), -1)
		r.writeIntToFSM(ctx, dialogManager, r.checkedItemKey(), selectedItem.ID)

		if r.onItemClick != nil {
			r.onItemClick(ctx, dialogManager, update, selectedItem)
		}

		if previousItemID != selectedItem.ID && r.onStateChanged != nil {
			r.onStateChanged(ctx, dialogManager, update, previousItemID, selectedItem.ID)
		}
	}

	return handlers
}

// callbackData формирует callback_data для конкретного элемента Radio.
func (r *Radio) callbackData(itemID int) string {
	return fmt.Sprintf("radio:%s:%d", r.ID, itemID)
}

// callbackPrefix возвращает префикс callback_data для маршрутизации обработчика Radio.
func (r *Radio) callbackPrefix() string {
	return fmt.Sprintf("radio:%s:", r.ID)
}

// checkedItemKey возвращает ключ состояния выбранного элемента в FSM.
func (r *Radio) checkedItemKey() string {
	return fmt.Sprintf("radio:%s:checked_item_id", r.ID)
}

// parseItemID извлекает ID выбранного элемента из callback update.
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

// findItemByID ищет элемент по ID в списке и возвращает его вместе с признаком успеха.
func (r *Radio) findItemByID(items []RadioItem, itemID int) (RadioItem, bool) {
	for _, item := range items {
		if item.ID == itemID {
			return item, true
		}
	}
	return RadioItem{}, false
}

// resolveHandlerData получает актуальные getter data для текущего окна через DialogManager.
func (r *Radio) resolveHandlerData(ctx context.Context, dialogManager *DialogManager) map[string]interface{} {
	if dialogManager == nil || dialogManager.dialog == nil || dialogManager.FSM == nil {
		return map[string]interface{}{}
	}

	window := dialogManager.dialog.GetWindow(*dialogManager.FSM.GetState())
	if window == nil {
		return map[string]interface{}{}
	}

	return window.getGetterData(ctx, dialogManager)
}

// renderItemText рендерит текст кнопки элемента по шаблону с учётом checked-состояния.
func (r *Radio) renderItemText(ctx context.Context, fsm FSMContext, data map[string]interface{}, item RadioItem) string {
	tmpl := r.uncheckedTextTemplate
	if r.IsChecked(ctx, fsm, item.ID) {
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
	templateData["checked"] = r.IsChecked(ctx, fsm, item.ID)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return ""
	}

	return buf.String()
}

// resolveItems достаёт список элементов Radio из getter data по ключу r.Items.
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

// readIntFromDialogManagerFSM читает числовое значение из FSM через DialogManager, либо возвращает fallback.
func (r *Radio) readIntFromDialogManagerFSM(ctx context.Context, dialogManager *DialogManager, key string, fallback int) int {
	if dialogManager == nil || dialogManager.FSM == nil {
		return fallback
	}
	return r.readIntFromFSM(*dialogManager.FSM, key, fallback)
}

// readIntFromFSM читает int по ключу из текущего state FSM, иначе возвращает fallback.
func (r *Radio) readIntFromFSM(fsm FSMContext, key string, fallback int) int {
	state := fsm.GetState()
	if state == nil {
		return fallback
	}

	value, exists := state.getData()[key]
	if !exists || value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

// writeIntToFSM сохраняет числовое значение в FSM как строку по переданному ключу.
func (r *Radio) writeIntToFSM(ctx context.Context, dialogManager *DialogManager, key string, value int) {
	if dialogManager == nil || dialogManager.FSM == nil {
		return
	}

	dialogManager.FSM.UpdateStateData(ctx, key, strconv.Itoa(value))
}
