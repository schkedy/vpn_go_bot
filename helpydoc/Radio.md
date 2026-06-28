## Radio
Зачем? — показывать список вариантов и давать пользователю выбрать один вариант кнопкой.

### Жизненный цикл

Рендерится окно -> показываются элементы кнопками -> пользователь нажимает кнопку -> приходит callback -> обновляется выбранный элемент -> выполняются колбэки -> окно рендерится заново.

---

## Текущая реализация

### Что такое Radio

`Radio` — stateful widget. Он хранит, какой элемент выбран сейчас (`checkedItemID`), и при новом нажатии снимает старый выбор и ставит новый.

### Структура

```go
Radio{
    ID             string
    Items          string
    Width          int
    checkedItemID  int
    onItemClick    RadioOnItemClick
    onStateChanged RadioOnStateChanged
    dialogManager  *DialogManager
}
```

`RadioItem` строго типизирован:

```go
type RadioItem struct {
    ID    int
    Value string
}
```

### Откуда берутся элементы

`Radio` получает элементы из `data[Items]` во время рендера.

Ожидаемый формат: только `[]RadioItem`.
Если в `data[Items]` лежит другой тип, кнопки не рендерятся.

Пример:

```go
map[string]interface{}{
    "fruits": []RadioItem{
        {ID: 1, Value: "Apple"},
        {ID: 2, Value: "Pear"},
        {ID: 3, Value: "Orange"},
    },
}
```

### Текст кнопок

У `Radio` два шаблона:

- `checkedText` (например, `"🔘 {{.item}}"`)
- `uncheckedText` (например, `"⚪️ {{.item}}"`)

В шаблоне доступны:

- `{{.item}}` — `Value` (`string`)
- `{{.item_id}}` — `ID` (`int`)
- `{{.checked}}` — `bool`
- остальные ключи из `data`

### Callback data

Формат кнопки:

```
radio:<RadioID>:<ItemID>
```

Пример: `radio:r_fruits:2`.

Обработчик у `Radio` регистрируется по префиксу `radio:<RadioID>:` и внутри парсит `ItemID` как `int`.

### Рендер

1. `Window` получает `data` из getter
2. `renderWidgets(data)` вызывает `radio.getButtonRows(data)`
3. `Radio` берёт `[]RadioItem` из `data[Items]`
4. для каждого элемента строит кнопку и callback data
5. кнопки раскладываются по строкам по `Width`

### Нажатие кнопки

1. Приходит callback `radio:<RadioID>:<ItemID>`
2. `Radio` парсит `ItemID` (`int`)
3. Берёт актуальные данные из `handlerData.GetterData`
4. Находит элемент, обновляет `checkedItemID`
5. Вызывает:
   - `onItemClick(...)`
   - `onStateChanged(...)` (только если выбор реально изменился)

### Управление выбором извне

```go
radio.GetChecked()    // int, -1 если ничего не выбрано
radio.IsChecked(2)    // bool
radio.SetChecked(3)   // принудительный выбор
```

### Ограничения

- `checkedItemID` хранится в памяти процесса.
- При перезапуске бота состояние сбрасывается.
- Для персистентности сохраняйте выбор в `DialogManager.Data` / Redis.
