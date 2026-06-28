package dialog

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type DialogManager struct {
	Session *DialogSession
	dialog  *Dialog
	sender  tgbotapi.BotAPI // TODO sender должен уметь отправлять сообщения, редактировать, удалять и т.д. в зависимости от того что нужно для рендера окна
	deps map[string]interface{}
}

// Dialog Manager Context uses for data which you get
// Хранится в Redis, ключ = userID (или chatID)
type DialogSession struct {
	State     State
	MessageID int
	ChatID    int64
	UserID    int64
	Data      map[string]interface{}
}

func NewDialogManager(dialog *Dialog) *DialogManager {
	return &DialogManager{
		dialog:  dialog,
		Session: &DialogSession{},
		sender:  Sender,
	}
}

// TODO: сделать проверку на то предыдущее сообщение либо текстовое либо медиа и
// в зависимости от этого редактировать его или отправлять новое
// TODO: добавить дату для сессии
func (dm *DialogManager) RenderWindow() {
	currentState := dm.Session.State
	if dm.dialog == nil {
		return
	}
	msgConfig := dm.dialog.GetWindow(currentState).RenderAll()
	if msgConfig.Media != nil {
		mediaMsg := dm.RenderMedia(msgConfig)
		if mediaMsg != nil {
			// dm.sender.DeleteMessage(dm.Session.ChatID,dm.Session.MessageID)
			// TODO: добавить SenderMode Edit или Send в зависимости от того нужно редактировать сообщение или отправлять новое
			sentMsg, err := dm.sender.EditMessageMediaAndMarkup(
				dm.Session.ChatID,
				dm.Session.MessageID,
				mediaMsg,
				msgConfig.Keyboard,
			)
			if err != nil {
				// TODO обработать ошибку
				return
			}
			dm.Session.MessageID = sentMsg.MessageID
			return
		}
	} else {
		msg := tgbotapi.NewMessage(dm.Session.ChatID, msgConfig.Text)
		msg.ReplyMarkup = msgConfig.Keyboard
		dm.sender.DeleteMessage(dm.Session.ChatID, dm.Session.MessageID)
		sentMsg, err := dm.sender.Send(msg)
		if err != nil {
			// TODO обработать ошибку
			return
		}
		dm.Session.MessageID = sentMsg.MessageID
	}

}

func (dm *DialogManager) RenderMedia(msgConfig WindowConfig) *tgbotapi.BaseInputMedia {
	switch msgConfig.Media.Type {
	case "photo":
		photoMsg := tgbotapi.NewBaseInputMedia("photo", buildFileData(msgConfig.Media))
		photoMsg.Caption = msgConfig.Text
		return &photoMsg
	case "video":
		videoMsg := tgbotapi.NewBaseInputMedia("video", buildFileData(msgConfig.Media))
		videoMsg.Caption = msgConfig.Text
		return &videoMsg
	case "audio":
		audioMsg := tgbotapi.NewBaseInputMedia("audio", buildFileData(msgConfig.Media))
		audioMsg.Caption = msgConfig.Text
		return &audioMsg
	case "document":
		docMsg := tgbotapi.NewBaseInputMedia("document", buildFileData(msgConfig.Media))
		docMsg.Caption = msgConfig.Text
		return &docMsg
	default:
		return nil
	}
}

func (dm *DialogManager) SwitchTo(state State) {
	// TODO: рендерим окно
	dm.Session.State = state

	if dm.dialog == nil {
		return
	}

	window := dm.dialog.GetWindow(state)
	if window != nil {
		window.BindDialogManager(dm)
	}
}

func (dm *DialogManager) Start(state State, data map[string]interface{}) {
	// TODO: рендерим окно
	dm.Session.State = state
	dm.Data = data

	if dm.dialog == nil {
		return
	}

	window := dm.dialog.GetWindow(state)
	if window != nil {
		window.BindDialogManager(dm)
	}
}
