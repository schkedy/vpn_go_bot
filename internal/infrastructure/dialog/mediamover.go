package dialog

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// MediaMoverOnSelect — пользовательский колбэк при выборе медиа пользователем.
type MediaMoverOnSelect func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update, index int, media Media)

// MediaMover — виджет, отображающий список медиа в виде пронумерованных кнопок с постраничной навигацией.
// Не хранит состояние внутри себя: текущая страница и выбранный индекс берутся из FSMContext.
type MediaMover struct {
	ID       string
	Items    string // ключ в getter data, по которому лежит []Media
	PageSize int    // количество кнопок-медиа на одной странице

	onSelect MediaMoverOnSelect
}

func NewMediaMover(id string, items string, pageSize int, onSelect MediaMoverOnSelect) *MediaMover {
	if pageSize <= 0 {
		pageSize = 1
	}

	return &MediaMover{
		ID:       id,
		Items:    items,
		PageSize: pageSize,
		onSelect: onSelect,
	}
}

// GetSelectedIndex возвращает глобальный индекс выбранного медиа или -1.
func (m *MediaMover) GetSelectedIndex(ctx context.Context, fsm FSMContext) int {
	_ = ctx
	return m.readIntFromFSM(fsm, m.selectedIndexKey(), -1)
}

// GetCurrentPage возвращает текущую страницу (0-based).
func (m *MediaMover) GetCurrentPage(ctx context.Context, fsm FSMContext) int {
	_ = ctx
	page := m.readIntFromFSM(fsm, m.currentPageKey(), 0)
	if page < 0 {
		return 0
	}
	return page
}

// SetSelectedIndex принудительно устанавливает выбранный индекс в FSMContext.
func (m *MediaMover) SetSelectedIndex(ctx context.Context, dialogManager *DialogManager, idx int) {
	m.writeIntToFSM(ctx, dialogManager, m.selectedIndexKey(), idx)
}

// GetCurrentMedia реализует MediaProvider — возвращает текущее выбранное медиа как *Media для рендера окна.
func (m *MediaMover) GetCurrentMedia(ctx context.Context, fsm FSMContext, data map[string]interface{}) *Media {
	items := m.resolveItems(data)
	idx := m.GetSelectedIndex(ctx, fsm)

	if idx < 0 || idx >= len(items) {
		return nil
	}

	copyItem := items[idx]
	return &copyItem
}

// getButtonRows строит строки кнопок:
//  1. Строка пронумерованных кнопок (1, 2, ...) для медиа текущей страницы
//  2. Строка пагинации (< · N/Total · >) — только если элементов больше чем PageSize
func (m *MediaMover) getButtonRows(ctx context.Context, fsm FSMContext, data map[string]interface{}) []ButtonRow {
	if m == nil {
		return nil
	}

	items := m.resolveItems(data)
	if len(items) == 0 {
		return nil
	}

	page := m.GetCurrentPage(ctx, fsm)
	totalPages := (len(items) + m.PageSize - 1) / m.PageSize
	if page >= totalPages {
		page = totalPages - 1
	}

	start := page * m.PageSize
	end := start + m.PageSize
	if end > len(items) {
		end = len(items)
	}

	pageItems := items[start:end]

	mediaRow := make([]*Button, 0, len(pageItems))
	for i := range pageItems {
		globalIdx := start + i
		label := fmt.Sprintf("%d", globalIdx+1)
		cb := m.callbackSelect(globalIdx)
		mediaRow = append(mediaRow, NewButton(label, cb, nil, ""))
	}

	rows := []ButtonRow{{Buttons: mediaRow}}

	if len(items) > m.PageSize {
		pageLabel := fmt.Sprintf("%d/%d", page+1, totalPages)
		prevBtn := NewButton("‹", m.callbackPrev(), nil, "")
		pageBtn := NewButton(pageLabel, m.callbackPageNoop(), nil, "")
		nextBtn := NewButton("›", m.callbackNext(), nil, "")

		rows = append(rows, ButtonRow{Buttons: []*Button{prevBtn, pageBtn, nextBtn}})
	}

	return rows
}

// getHandlers регистрирует обработчики для всех кнопок медиа + пагинации.
func (m *MediaMover) getHandlers() map[string]HandlerFunc {
	handlers := make(map[string]HandlerFunc)
	if m == nil {
		return handlers
	}

	handlers[m.callbackSelectPrefix()] = func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update) {
		idx := m.parseSelectedIndex(update)
		if idx < 0 {
			return
		}

		items := m.resolveItems(m.resolveHandlerData(ctx, dialogManager))
		if idx >= len(items) {
			return
		}

		m.writeIntToFSM(ctx, dialogManager, m.selectedIndexKey(), idx)
		if m.PageSize > 0 {
			m.writeIntToFSM(ctx, dialogManager, m.currentPageKey(), idx/m.PageSize)
		}

		if m.onSelect != nil {
			m.onSelect(ctx, dialogManager, update, idx, items[idx])
		}
	}

	handlers[m.callbackPrev()] = func(ctx context.Context, dialogManager *DialogManager, _ *tgbotapi.Update) {
		items := m.resolveItems(m.resolveHandlerData(ctx, dialogManager))
		totalPages := (len(items) + m.PageSize - 1) / m.PageSize
		if totalPages == 0 {
			return
		}

		page := m.readIntFromDialogManagerFSM(ctx, dialogManager, m.currentPageKey(), 0)
		if page > 0 {
			page--
		} else {
			page = totalPages - 1
		}

		m.writeIntToFSM(ctx, dialogManager, m.currentPageKey(), page)
	}

	handlers[m.callbackNext()] = func(ctx context.Context, dialogManager *DialogManager, _ *tgbotapi.Update) {
		items := m.resolveItems(m.resolveHandlerData(ctx, dialogManager))
		totalPages := (len(items) + m.PageSize - 1) / m.PageSize
		if totalPages == 0 {
			return
		}

		page := m.readIntFromDialogManagerFSM(ctx, dialogManager, m.currentPageKey(), 0)
		if page < totalPages-1 {
			page++
		} else {
			page = 0
		}

		m.writeIntToFSM(ctx, dialogManager, m.currentPageKey(), page)
	}

	handlers[m.callbackPageNoop()] = func(_ context.Context, _ *DialogManager, _ *tgbotapi.Update) {}

	return handlers
}

func (m *MediaMover) callbackSelect(idx int) string {
	return fmt.Sprintf("mediamover:%s:select:%d", m.ID, idx)
}

func (m *MediaMover) callbackSelectPrefix() string {
	return fmt.Sprintf("mediamover:%s:select:", m.ID)
}

func (m *MediaMover) callbackPrev() string {
	return fmt.Sprintf("mediamover:%s:prev", m.ID)
}

func (m *MediaMover) callbackNext() string {
	return fmt.Sprintf("mediamover:%s:next", m.ID)
}

func (m *MediaMover) callbackPageNoop() string {
	return fmt.Sprintf("mediamover:%s:page", m.ID)
}

func (m *MediaMover) selectedIndexKey() string {
	return fmt.Sprintf("mediamover:%s:selected_idx", m.ID)
}

func (m *MediaMover) currentPageKey() string {
	return fmt.Sprintf("mediamover:%s:current_page", m.ID)
}

func (m *MediaMover) parseSelectedIndex(update *tgbotapi.Update) int {
	if update == nil || update.CallbackQuery == nil {
		return -1
	}

	data := update.CallbackQuery.Data
	prefix := m.callbackSelectPrefix()
	if len(data) <= len(prefix) || data[:len(prefix)] != prefix {
		return -1
	}

	idx, err := strconv.Atoi(data[len(prefix):])
	if err != nil {
		return -1
	}

	return idx
}

func (m *MediaMover) resolveHandlerData(ctx context.Context, dialogManager *DialogManager) map[string]interface{} {
	if dialogManager == nil || dialogManager.dialog == nil || dialogManager.FSM == nil {
		return map[string]interface{}{}
	}

	window := dialogManager.dialog.GetWindow(*dialogManager.FSM.GetState())
	if window == nil {
		return map[string]interface{}{}
	}

	return window.getGetterData(ctx, dialogManager)
}

// resolveItems возвращает слайс медиа из getter data по ключу m.Items.
func (m *MediaMover) resolveItems(data map[string]interface{}) []Media {
	if m == nil || m.Items == "" || data == nil {
		return nil
	}

	raw, exists := data[m.Items]
	if !exists || raw == nil {
		return nil
	}

	items, ok := raw.([]Media)
	if !ok {
		return nil
	}

	return items
}

func (m *MediaMover) readIntFromDialogManagerFSM(ctx context.Context, dialogManager *DialogManager, key string, fallback int) int {
	if dialogManager == nil || dialogManager.FSM == nil {
		return fallback
	}
	return m.readIntFromFSM(*dialogManager.FSM, key, fallback)
}

func (m *MediaMover) readIntFromFSM(fsm FSMContext, key string, fallback int) int {
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

func (m *MediaMover) writeIntToFSM(ctx context.Context, dialogManager *DialogManager, key string, value int) {
	if dialogManager == nil || dialogManager.FSM == nil {
		return
	}

	dialogManager.FSM.UpdateStateData(ctx, key, strconv.Itoa(value))
}
