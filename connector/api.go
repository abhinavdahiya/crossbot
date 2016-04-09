package connector

import "net/url"

type BotAPI interface {
	SendMessage(User, string, Options) error
	SendPhoto(User, PhotoSize, Options) error
	SendAudio(User, Audio, Options) error
	SendDocument(User, Document, Options) error
	SendVideo(User, Video, Options) error
	GetFileURL(string) string

	GetUpdates(chan<- Message, chan<- struct{}, chan<- error)
	SetAndListen(string, string, chan<- Message, chan<- struct{}, chan<- error)
}

type Options struct {
	ParseMode             string
	Caption               string
	DisableWebPagePreview bool
	Keyboard              [][]string
	ResizeKeyboard        bool
	OneTimeKeyboard       bool
	HideKeyboard          bool
}

// UpdateConfig contains information about a GetUpdates request.
type UpdateConfig struct {
	Offset  int64
	Limit   int64
	Timeout int64
}

// WebhookConfig contains information about a SetWebhook request.
type WebhookConfig struct {
	URL         *url.URL
	Certificate interface{}
}

// NewUpdate gets updates since the last Offset.
//
// offset is the last Update ID to include.
// You likely want to set this to the last Update ID plus 1.
func NewUpdate(offset int64) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}
