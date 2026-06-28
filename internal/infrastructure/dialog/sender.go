package dialog

type Sender interface {
	SendMessage(chatID int64, text string) (messageID int, err error)
	EditMessage(chatID int64, messageID int, newText string) error
	DeleteMessage(chatID int64, messageID int) error
}
