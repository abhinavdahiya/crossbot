package crossbot

import "strings"

// General Message received from connector
type Message struct {
	MessageID int64
	From      User
	Date      int64
	Chat      Chat
	Bot       *BotAPI
}

// This struct is used to identify a user
type User struct {
	ID        int64
	FirstName string
	LastName  string
	Handle    string
	IsGroup   bool
}

// This stores the contents received from bot
type Chat struct {
	Type     string
	Text     string
	Photo    []PhotoSize
	Audio    Audio
	Video    Video
	Document Document
	Voice    Voice
	Location Location
}

// Checks if the message is a command
func (c Chat) IsCommand() bool {
	return c.Type == "command"
}

// This func return the command received in the message
// ex: /command => command
// ex: /command@bot => command
func (c Chat) Command() string {
	if !c.IsCommand() {
		return ""
	}

	command := strings.SplitN(c.Text, " ", 2)[0][1:]

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// PhotoSize contains information about photos.
type PhotoSize struct {
	FileID   string
	Width    int
	Height   int
	FileSize int // optional
}

// Audio contains information about audio.
type Audio struct {
	FileID    string
	Duration  int
	Performer string // optional
	Title     string // optional
	MimeType  string // optional
	FileSize  int    // optional
}

// Document contains information about a document.
type Document struct {
	FileID    string
	Thumbnail PhotoSize // optional
	FileName  string    // optional
	MimeType  string    // optional
	FileSize  int       // optional
}

// Video contains information about a video.
type Video struct {
	FileID    string
	Width     int
	Height    int
	Duration  int
	Thumbnail PhotoSize // optional
	MimeType  string    // optional
	FileSize  int       // optional
}

// Voice contains information about a voice.
type Voice struct {
	FileID   string
	Duration int
	MimeType string // optional
	FileSize int    // optional
}

// Location contains information about a place.
type Location struct {
	Longitude float32
	Latitude  float32
}
