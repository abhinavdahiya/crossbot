package crossbot

import (
	"sync"
	"time"

	"github.com/abhinavdahiya/crossbot/connector"
)

var (
	sMutex sync.RWMutex
)

const ()

type Store struct {
	ID    int64
	State State
	time  int64
}

type Bot struct {
	States    map[string]State
	Storage   map[connector.User]Store
	InitState string
}

func NewBot() *Bot {
	return &Bot{
		States:    make(map[string]State),
		Storage:   make(map[connector.User]Store),
		InitState: "start",
	}
}

func (b *Bot) AddState(s State) {
	b.States[s.Name] = s
}

func (b *Bot) StoreState(u connector.User, mID int64, s State) {
	sMutex.Lock()
	b.Storage[u] = Store{
		ID:    mID,
		State: s,
		time:  time.Now().Unix(),
	}
	sMutex.Unlock()
}

func (b *Bot) FetchState(u connector.User) State {
	sMutex.RLock()
	if s, ok := b.Storage[u]; ok {
		sMutex.RUnlock()
		return s.State
	}
	sMutex.RUnlock()
	return Store{}
}

// The message is processed as follows:
// Load the current state of User
// Exec Leave action on State
// Test for next state if IsMoved == false
// Transfer ctx from current to next state
// Exec Enter of next state
// Store next state for user
func (b *Bot) Process(m connector.Message) error {
	curr := b.FetchState(m.From)
	if curr == (Store{}) {
		// set user state to InitState
		tmp := b.States[b.InitState]
		b.StoreState(m.From, m.MessageID, tmp)
		curr = tmp
	}

	if curr.Leave != nil {
		err := curr.Leave(m)
		if err != nil {
			return err
		}
	}

	var next string
	if curr.IsMoved {
		next = curr.Chain
	}

	var err error
	next, err = curr.Test(m)
	if err != nil {
		return err
	}

	ctx := curr.Data()
	next.SetData(ctx)

	err = next.Enter(m)
	if err != nil {
		return err
	}

	b.StoreState(m.From, m.MessageID, next)
}
