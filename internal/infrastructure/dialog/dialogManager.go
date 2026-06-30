package dialog

import (
	"context"
	"errors"
	"strconv"
	"vpn_go_bot/internal/infrastructure/cache"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

var NoDialogSessionInStorageError error = errors.New("no dialog session in storage")

type DialogManager struct {
	Session *DialogSession
	dialog  *Dialog
	sender  tgbotapi.BotAPI // TODO sender должен уметь отправлять сообщения, редактировать, удалять и т.д. в зависимости от того что нужно для рендера окна
	FSM     *FSMContext
	deps    map[string]interface{}
}

//

// Dialog Manager Context uses for data which you get
// Хранится в Redis, ключ = userID (или chatID)
type DialogSession struct {
	MessageID int
	ChatID    int64
	UserID    int64
	Data      map[string]string
}

// TODO : Сделать приведение string к interface{}
func newDialogSession(userID int64, chatID int64, messageID int) *DialogSession {
	return &DialogSession{
		MessageID: messageID,
		ChatID:    chatID,
		UserID:    userID,
		Data:      make(map[string]string),
	}
}

func NewDialogSessionFromStorage(ctx context.Context, storage *cache.RedisClient, userID int64) (*DialogSession, error) {
	hashKeySession := "dialog_session:" + string(userID)
	hashKeyData := "dialog_data:" + string(userID)
	sessionData, err := storage.HGetAll(ctx, hashKeySession)
	if sessionData == nil {
		return nil, NoDialogSessionInStorageError
	}
	if err != nil {
		return nil, err
	}
	data, err := storage.HGetAll(ctx, hashKeyData)
	if err != nil {
		return nil, err
	}
	MessageID, err := strconv.Atoi(sessionData["MessageID"])
	if err != nil {
		return nil, err
	}
	ChatID, err := strconv.ParseInt(sessionData["ChatID"], 10, 64)
	if err != nil {
		return nil, err
	}
	UserID, err := strconv.ParseInt(sessionData["UserID"], 10, 64)
	if err != nil {
		return nil, err
	}

	dialogSession := &DialogSession{
		MessageID: MessageID,
		ChatID:    ChatID,
		UserID:    UserID,
		Data:      data,
	}
	return dialogSession, nil

}

func (ds *DialogSession) SaveToStorage(ctx context.Context, storage *cache.RedisClient) error {
	hashKeySession := "dialog_session:" + string(ds.UserID)
	hashKeyData := "dialog_data:" + string(ds.UserID)
	err := storage.HSet(ctx, hashKeySession, "MessageID", ds.MessageID)
	if err != nil {
		return err
	}
	err = storage.HSet(ctx, hashKeySession, "ChatID", ds.ChatID)
	if err != nil {
		return err
	}
	err = storage.HSet(ctx, hashKeySession, "UserID", ds.UserID)
	if err != nil {
		return err
	}
	for key, value := range ds.Data {
		err = storage.HSet(ctx, hashKeyData, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO : закончить, нужно правильно сохранениеи DialogSession.Data
// update,
func NewDialogManager(
	ctx context.Context,
	update tgbotapi.Update,
	dialog *Dialog,
	sender tgbotapi.BotAPI,
	FSM *FSMContext,
	deps map[string]interface{},
	storage *cache.RedisClient,
) (*DialogManager, error) {
	dialogSession, err := NewDialogSessionFromStorage(ctx, storage, int64(update.Message.From.ID))
	switch err {
	case NoDialogSessionInStorageError:
		dialogSession = newDialogSession(int64(update.Message.From.ID), update.Message.Chat.ID, -1)
	case nil:
	default:
		return nil, err
	}

	dialogManager := &DialogManager{
		Session: dialogSession,
		dialog:  dialog,
		sender:  sender,
		FSM:     FSM,
		deps:    deps,
	}
	return dialogManager, nil
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
