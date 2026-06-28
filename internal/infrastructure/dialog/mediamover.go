package dialog

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// TODO MediaMover.ID  должен браться на  основе прикрпленния к пользователю диалога или нет, надо узнать как колбэк будет понимать
// к какому именно MediaMover он относится, если таких виджетов несколько в одном окне

// MediaMoverOnSelect — пользовательский колбэк при выборе медиа пользователем.
type MediaMoverOnSelect func(ctx context.Context, dialogManager *DialogManager, update *tgbotapi.Update, index int, media Media)

// MediaMover — виджет, отображающий список медиа в виде пронумерованных кнопок с постраничной навигацией.
// Реализует Widget и MediaProvider.
type MediaMover struct {
	ID       string
	Items    string // ключ в getter data, по которому лежит []Media
	PageSize int    // количество кнопок-медиа на одной странице

	onSelect MediaMoverOnSelect

	mu          sync.RWMutex
	currentPage int // 0-based
	selectedIdx int // глобальный индекс выбранного медиа, -1 = ничего не выбрано
}

func NewMediaMover(id string, items string, pageSize int, onSelect MediaMoverOnSelect) *MediaMover {
	if pageSize <= 0 {
		pageSize = 1
	}

	return &MediaMover{
		ID:          id,
		Items:       items,
		PageSize:    pageSize,
		onSelect:    onSelect,
		selectedIdx: -1,
	}
}

// GetSelectedIndex возвращает глобальный индекс выбранного медиа или -1.
func (m *MediaMover) GetSelectedIndex() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.selectedIdx
}

// GetCurrentPage возвращает текущую страницу (0-based).
func (m *MediaMover) GetCurrentPage() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.currentPage
}

// SetSelectedIndex принудительно устанавливает выбранный индекс.
func (m *MediaMover) SetSelectedIndex(idx int) {
	m.mu.Lock()
	m.selectedIdx = idx
	m.mu.Unlock()
}

// GetCurrentMedia реализует MediaProvider — возвращает текущее выбранное медиа как *Media для рендера окна.
func (m *MediaMover) GetCurrentMedia(data map[string]interface{}) *Media {
	items := m.resolveItems(data)

	m.mu.RLock()
	idx := m.selectedIdx
	m.mu.RUnlock()

	if idx < 0 || idx >= len(items) {
		return nil
	}

	copy := items[idx]
	return &copy
}

// getButtonRows строит строки кнопок:
//  1. Строка пронумерованных кнопок (1, 2, ...) для медиа текущей страницы
//  2. Строка пагинации (< · N/Total · >) — только если элементов больше чем PageSize
func (m *MediaMover) getButtonRows(data map[string]interface{}) []ButtonRow {
	if m == nil {
		return nil
	}

	items := m.resolveItems(data)
	if len(items) == 0 {
		return nil
	}

	m.mu.RLock()
	page := m.currentPage
	m.mu.RUnlock()

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

	// строка с пронумерованными кнопками медиа
	mediaRow := make([]*Button, 0, len(pageItems))
	for i, item := range pageItems {
		globalIdx := start + i
		label := fmt.Sprintf("%d", globalIdx+1)
		cb := m.callbackSelect(globalIdx)

		_ = item // globalIdx несёт всю нужную информацию; item используется в handler
		btn := NewButton(label, cb, nil, "")
		mediaRow = append(mediaRow, btn)
	}

	rows := []ButtonRow{{Buttons: mediaRow}}

	// строка пагинации — только если элементов больше чем PageSize
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

		items := m.resolveItems(m.resolveHandlerData(handlerData))
		if idx >= len(items) {
			return
		}

		m.mu.Lock()
		m.selectedIdx = idx
		m.mu.Unlock()

		if m.onSelect != nil {
			m.onSelect(ctx, handlerData, update, idx, items[idx])
		}
	}

	handlers[m.callbackPrev()] = func(_ context.Context, handlerData *handler.HandlerData, update *tgbotapi.Update) {
		items := m.resolveItems(m.resolveHandlerData(handlerData))
		totalPages := (len(items) + m.PageSize - 1) / m.PageSize
		if totalPages == 0 {
			return
		}

		m.mu.Lock()
		if m.currentPage > 0 {
			m.currentPage--
		} else {
			m.currentPage = totalPages - 1
		}
		m.mu.Unlock()
	}

	handlers[m.callbackNext()] = func(_ context.Context, handlerData *handler.HandlerData, update *tgbotapi.Update) {
		items := m.resolveItems(m.resolveHandlerData(handlerData))
		totalPages := (len(items) + m.PageSize - 1) / m.PageSize
		if totalPages == 0 {
			return
		}

		m.mu.Lock()
		if m.currentPage < totalPages-1 {
			m.currentPage++
		} else {
			m.currentPage = 0
		}
		m.mu.Unlock()
	}

	handlers[m.callbackPageNoop()] = func(_ context.Context, _ *handler.HandlerData, _ *tgbotapi.Update) {}

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

func (m *MediaMover) parseSelectedIndex(update *tgbotapi.Update) int {
	if update == nil || update.CallbackQuery == nil {
		return -1
	}

	data := update.CallbackQuery.Data
	prefix := m.callbackSelectPrefix()
	if len(data) <= len(prefix) || data[:len(prefix)] != prefix {
		return -1
	}

	var idx int
	if _, err := fmt.Sscanf(data[len(prefix):], "%d", &idx); err != nil {
		return -1
	}

	return idx
}

func (m *MediaMover) resolveHandlerData(handlerData *handler.HandlerData) map[string]interface{} {
	if handlerData == nil || handlerData.GetterData == nil {
		return map[string]interface{}{}
	}

	return handlerData.GetterData
}

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
