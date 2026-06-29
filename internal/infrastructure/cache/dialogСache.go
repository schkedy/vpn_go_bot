package cache

import (
    "context"
    "errors"
)

var (
    ErrSessionNotFound     = errors.New("dialog session not found")
    ErrWidgetStateNotFound = errors.New("widget state not found")
)

// DialogSessionSnapshot — минимальный сериализуемый срез сессии диалога.
type DialogSessionSnapshot struct {
    MessageID  int64                  `json:"message_id"`
    State      string                 `json:"state"`
    DialogData map[string]interface{} `json:"dialog_data"`
}

// DialogStateStore — единый контракт для хранения состояния диалога и виджетов.
type DialogStateStore interface {
    GetSession(ctx context.Context, userID int64) (*DialogSessionSnapshot, error)
    SaveSession(ctx context.Context, userID int64, session *DialogSessionSnapshot) error
    DeleteSession(ctx context.Context, userID int64) error

    GetWidgetState(ctx context.Context, userID int64, state string, widgetID string, out interface{}) error
    SetWidgetState(ctx context.Context, userID int64, state string, widgetID string, in interface{}) error
    DeleteWidgetState(ctx context.Context, userID int64, state string, widgetID string) error
}
