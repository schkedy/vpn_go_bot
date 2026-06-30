package dialog

import (
	"context"
	"errors"
	"strconv"
	"vpn_go_bot/internal/infrastructure/cache"
)

var (
	SessionKeyPrefix                    = "dialog_session"
	DataKeyPrefix                       = "dialog_data"
	NoDialogSessionInStorageError error = errors.New("no dialog session in storage")
)

// TODO: заменить redis на интерфейс, чтобы можно было использовать разные хранилища
type DialogSessionStorage struct {
	redis *cache.RedisClient
}

func NewDialogSessionStorage(redis *cache.RedisClient) *DialogSessionStorage {
	return &DialogSessionStorage{
		redis: redis,
	}
}

func (dss *DialogSessionStorage) SaveSession(ctx context.Context, session *DialogSession) error {
	sessionKey := SessionKeyPrefix + ":" + string(session.UserID)
	dataKey := DataKeyPrefix + ":" + string(session.UserID)

	// Сохраняем данные сессии в Redis
	err := dss.redis.HSet(ctx, sessionKey, "MessageID", session.MessageID)
	if err != nil {
		return err
	}
	err = dss.redis.HSet(ctx, sessionKey, "ChatID", session.ChatID)
	if err != nil {
		return err
	}
	err = dss.redis.HSet(ctx, sessionKey, "UserID", session.UserID)
	if err != nil {
		return err
	}

	// Сохраняем данные диалога в Redis
	for key, value := range session.Data {
		err = dss.redis.HSet(ctx, dataKey, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dss *DialogSessionStorage) GetSession(ctx context.Context, userID int64) (*DialogSession, error) {
	SessionKey := SessionKeyPrefix + ":" + string(userID)
	DataKey := DataKeyPrefix + ":" + string(userID)

	// Получаем данные сессии из Redis
	sessionData, err := dss.redis.HGetAll(ctx, SessionKey)
	// TODO : проверить что действительно при остутивии в редис возвращается nil, а не пустая map
	if sessionData == nil {
		return nil, NoDialogSessionInStorageError
	}
	if err != nil {
		return nil, err
	}

	// Получаем данные диалога из Redis
	data, err := dss.redis.HGetAll(ctx, DataKey)
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
