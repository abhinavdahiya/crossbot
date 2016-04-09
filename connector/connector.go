package connector

import (
	"errors"
	"net/http"
	"os"
	"strconv"
)

const (
	ErrNoCert    = errors.New("Couldn't read certificate file or no certificate provided")
	ErrNoKey     = errors.New("Couldn't read key file or no key provided")
	ErrNoBaseURL = errors.New("No BaseURL found")
)

type Connector struct {
	Endpoints []*BotAPI

	CertFile string
	KeyFile  string
	BaseURL  string

	MsgCh  chan Message
	StopCh chan struct{}
}

// Compulsory SSL when setting webhooks
// This should be either empty or all fields must exist
type Config struct {
	CertFile string
	KeyFile  string
	BaseURL  string
}

func NewConnector(c Config) (*Connector, error) {
	conn := &Connector{
		MsgCh:  make(chan Message, 100),
		StopCh: make(chan struct{}),
	}

	if c != (Config{}) {
		if _, err := os.Stat(c.CertFile); os.IsNotExist(err) {
			return nil, ErrNoCert
		}

		if _, err := os.Stat(c.KeyFile); os.IsNotExist(err) {
			return nil, ErrNoKey
		}

		if url := c.BaseURL; url != "" {
			conn.CertFile = c.CertFile
			conn.KeyFile = c.KeyFile
			conn.BaseURL = url

			return conn, nil
		}

		return nil, ErrNoBaseURL
	}
}

func (c *Connector) AddEndpoint(b *BotAPI) {
	c.Endpoints = append(c.Endpoints, b)
}

func (c *Connector) Start() (<-chan Message, <-chan error) {
	ErrCh := make(chan error, 100)
	for _, eps := range c.Endpoints {
		eps.GetUpdates(c.MsgCh, c.Stop, ErrCh)
	}
	return c.MsgCh, ErrCh
}

func (c *Connector) StartWebhook(port int, msg chan<- Message, err chan<- error) {
	for _, eps := range c.Endpoints {
		eps.SetAndListen(c.BaseURL, c.CertFile, msg, c.StopCh, err)
	}
	p := strconv.Itoa(port)
	http.ListenAndServeTLS("0.0.0.0:"+p, c.CertFile, c.KeyFile, nil)
}

func (c *Connector) Stop() {
	close(c.StopCh)
}
