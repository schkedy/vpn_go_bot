## MediaMover
Зачем? — показывать пользователю список медиа как кнопки с пагинацией, позволять выбрать элемент и отрендерить выбранное медиа.

### Жизненный цикл

Рендерится окно -> MediaMover строит кнопки из данных -> пользователь нажимает номер/стрелку -> обновляется внутреннее состояние -> окно рендерится заново -> отображается актуальное медиа.

---

## Текущая реализация

### Что такое MediaMover

`MediaMover` — stateful widget с двумя состояниями:

- `selectedIdx` — выбранный глобальный индекс (`-1`, если ничего не выбрано)
- `currentPage` — текущая страница (0-based)

Реализует:

- `Widget` (кнопки + callback handlers)
- `MediaProvider` (возврат текущего медиа)

### Структура

```go
MediaMover{
    ID          string
    Items       string
    PageSize    int
    onSelect    MediaMoverOnSelect
    currentPage int
    selectedIdx int
}
```

### Откуда берутся элементы

На рендере берётся `data[Items]`, ожидается только `[]Media`.
Если тип другой, MediaMover ничего не рендерит.

### Формат callback data

Выбор медиа:

```
mediamover:<ID>:select:<globalIndex>
```

Пагинация:

```
mediamover:<ID>:prev
mediamover:<ID>:next
mediamover:<ID>:page
```

### Рендер кнопок

1. `Window` получает `data` из getter
2. `renderWidgets(data)` вызывает `mediaMover.getButtonRows(data)`
3. строятся кнопки номеров для текущей страницы
4. если элементов больше `PageSize`, добавляется строка пагинации

### Обработка callback

В обработчиках MediaMover использует не getter-функцию, а `handlerData.GetterData`:

- `select` — парсит индекс, обновляет `selectedIdx`, вызывает `onSelect`
- `prev` / `next` — меняет `currentPage` с циклическим переходом
- `page` — noop

### Как отдается медиа в окно

`Window.renderMedia(data)` сначала спрашивает все widget, реализующие `MediaProvider`.
Если `MediaMover.GetCurrentMedia(data)` вернул медиа, оно имеет приоритет над статическим медиа окна.

### Управление состоянием извне

```go
mover.GetSelectedIndex()  // int
mover.GetCurrentPage()    // int
mover.SetSelectedIndex(2)
```

### Ограничения

- `selectedIdx` и `currentPage` в памяти процесса.
- После перезапуска бота состояние теряется.
- Для персистентности сохраняйте стейт во внешнем хранилище.
