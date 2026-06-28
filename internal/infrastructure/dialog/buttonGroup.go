package dialog

// ButtonGroup объединяет несколько кнопок и раскладывает их по строкам согласно Width.
type ButtonGroup struct {
	Buttons []*Button
	Width   int
}

// getButtonRows раскладывает кнопки по строкам согласно Width и фильтрует их по ShowWhenKey.
func (g *ButtonGroup) getButtonRows(data map[string]interface{}) []ButtonRow {
	if g == nil || len(g.Buttons) == 0 {
		return nil
	}

	width := g.Width
	if width <= 0 {
		width = 1
	}

	visibleButtons := make([]*Button, 0, len(g.Buttons))
	for _, button := range g.Buttons {
		if button != nil && button.ShouldShow(data) {
			visibleButtons = append(visibleButtons, button)
		}
	}

	if len(visibleButtons) == 0 {
		return nil
	}

	rows := make([]ButtonRow, 0, (len(visibleButtons)+width-1)/width)
	for start := 0; start < len(visibleButtons); start += width {
		end := start + width
		if end > len(visibleButtons) {
			end = len(visibleButtons)
		}

		rows = append(rows, ButtonRow{Buttons: visibleButtons[start:end]})
	}

	return rows
}

// getHandlers return all binding HandlerFunc for buttons in group, key is button.CallbackData
func (g *ButtonGroup) getHandlers() map[string]HandlerFunc {
	handlers := make(map[string]HandlerFunc)
	if g == nil {
		return handlers
	}

	for _, button := range g.Buttons {
		for key, handlerFunc := range button.getHandlers() {
			handlers[key] = handlerFunc
		}
	}

	return handlers
}
