package cache

import (
	"context"
	"errors"
)

var ErrSessionNotFound = errors.New("dialog session not found")

type DialogSession struct {
	ChatID     int64                  `json:"chat_id"`
	UserID     int64                  `json:"user_id"`
	MessageID  int                    `json:"message_id"`
	State      string                 `json:"state"`
	DialogData map[string]interface{} `json:"dialog_data"`
}

type DialogCache interface {
	SaveSession(ctx context.Context, session *DialogSession) error
	GetSession(ctx context.Context, userID int64) (*DialogSession, error)
	DeleteSession(ctx context.Context, userID int64) error
}
