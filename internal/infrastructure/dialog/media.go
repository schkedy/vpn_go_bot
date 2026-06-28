package dialog

// Media — универсальная структура медиа-элемента.
// Покрывает оба сценария:
//   - Статическое медиа окна: заполняется при создании Window, HasMediaKey контролирует
//     условный показ через данные getter.
//   - Динамические элементы списка (MediaMover): заполняется из getter по ключу Items,
//     HasMediaKey не используется.
import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Media struct {
	Type        string // тип медиа: photo, video, audio, document
	FileID      string // Telegram file_id (для уже отправленных файлов)
	FilePath    string // путь до локального файла
	URL         string // URL на файл
	HasMediaKey string // ключ в getter для условного показа статического медиа (опционально)
}

// buildFileData make RequestFileData from FilePath any nature
func buildFileData(media *Media) tgbotapi.RequestFileData {
	if media == nil {
		return nil
	}

	if media.FileID != "" {
		return tgbotapi.FileID(media.FileID)
	}

	if media.FilePath != "" {
		return tgbotapi.FilePath(media.FilePath)
	}

	if media.URL != "" {
		return tgbotapi.FileURL(media.URL)
	}

	return nil
}
