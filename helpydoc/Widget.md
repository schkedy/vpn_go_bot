## Widget
Зачем нужен — чтобы описывать UI-элементы окна как независимые блоки (кнопки, группы, radio, media mover) и централизованно собирать их в `Window`.

---

## Текущая реализация

### Интерфейс Widget

Каждый widget должен уметь:

1. вернуть строки кнопок для рендера
2. вернуть map callback handlers

Концептуально:

- widget -> `[]ButtonRow`
- widget -> `map[callbackData]handler`

`Window` не знает внутреннюю реализацию widget — только объединяет их результаты.

### MediaProvider

Дополнительный интерфейс для widget, которые умеют отдавать текущее медиа:

- `GetCurrentMedia(data map[string]interface{}) *Media`

`Window.renderMedia(data)` сначала проверяет такие widget.

### Button

`Button` содержит:

1. шаблон текста
2. `CallbackData`
3. `Handler`
4. `ShowWhenKey`

Текст и видимость вычисляются по переданным `data`.

### Group

`Group` раскладывает кнопки по строкам:

- `Buttons []*Button`
- `Width int`

Пример: при `Width = 2` и 5 кнопках получится 3 строки (2 + 2 + 1).

### Роль Window

`Window`:

1. один раз вызывает getter и получает `data`
2. передаёт эти `data` в рендер widget
3. собирает итоговый `InlineKeyboardMarkup`
4. объединяет handlers от всех widget

Важно: widget больше не хранят getter-функцию и не вызывают её напрямую.

### Callback-путь данных

Для обработки callback актуальные данные передаются через `handler.HandlerData.GetterData`.
Это делает поток данных явным и убирает скрытые вызовы getter внутри widget.
