package crossbot

import "net/url"

type BotAPI interface {
	SendMessage(User, string, Options) error
	SendPhoto(User, Photo, Options) error
	SendVideo(User, Photo, Options) error
	GetFileURL(string) string

	GetUpdates(UpdateConfig) ([]Message, error)
	SetWebhook(WebhookConfig) (bool, error)
}

type Options struct {
	ParseMode             string
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

// NewWebhook creates a new webhook.
//
// link is the url parsable link you wish to get the updates.
func NewWebhook(link string) WebhookConfig {
	u, _ := url.Parse(link)

	return WebhookConfig{
		URL: u,
	}
}

// NewWebhookWithCert creates a new webhook with a certificate.
//
// link is the url you wish to get webhooks,
// file contains a string to a file, FileReader, or FileBytes.
func NewWebhookWithCert(link string, file interface{}) WebhookConfig {
	u, _ := url.Parse(link)

	return WebhookConfig{
		URL:         u,
		Certificate: file,
	}
}
