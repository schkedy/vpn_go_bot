# EditMessage TODO

Источник: `github.com/OvyFlash/telegram-bot-api` (`NewEditMessage...` constructors).

## Constructors найдено

- [x] `NewEditMessageMedia`
- [x] `NewEditMessagePhoto`
- [x] `NewEditMessageVideo`
- [x] `NewEditMessageAnimation`
- [x] `NewEditMessageAudio`
- [x] `NewEditMessageDocument`
- [x] `NewEditMessageText`
- [x] `NewEditMessageTextAndMarkup`
- [x] `NewEditMessageCaption`
- [x] `NewEditMessageReplyMarkup`
- [x] `NewEditMessageChecklist`

## BotAPI methods

- [x] `EditMessageMedia`
- [x] `EditMessagePhoto`
- [x] `EditMessageVideo`
- [x] `EditMessageAnimation`
- [x] `EditMessageAudio`
- [x] `EditMessageDocument`
- [x] `EditMessageText`
- [x] `EditMessageTextAndMarkup`
- [x] `EditMessageCaption`
- [x] `EditMessageReplyMarkup`
- [x] `EditMessageChecklist`

## Notes

- Нельзя менять сообщение с media на сообщение только с текстом, нужно отпрвалять его просто заново
- Удален обрыв объявления `func (bot *BotAPI) EditMessage`, который ломал компиляцию.
- Проверка: `go test ./...` в `/home/me/go/telegram-bot-api` — успешно.
