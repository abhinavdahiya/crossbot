package builder

import (
	"errors"

	"github.com/abhinavdahiya/crossbot/crossbot"
)

const (
	ErrNoRuleFound = errors.New("No Rule found to transition the state")
)

// This struct defines any state for the bot
// corresponding to any user
// This also stores Transistors for all possible
// moves from this state
//
// make sure the state Name is unique
type State struct {
	Name         string
	Enter, Leave Action
	IsMoved      bool
	Chain        string
	CmdRules     map[string]Transistor
	TypeRules    map[string]Transistor
	FallbackRule Transistor
}

// This gets the context data for a state
func (s *State) Data() interface{} {
	d := Get(s)
	return d
}

// Stores context for a state
func (s *State) SetData(d interface{}) {
	Set(s, d)
}

// This function forces the state change to new state
// and bypasses the Test func
func (s *State) Transit(ns string) {
	s.IsMoved = true
	s.Chain = ns
}

// Registers a command transistor for this state
func (s *State) RegisterCmdRule(cmd string, t Transistor) {
	s.CmdRules[cmd] = t
}

// Registers a transistor corresponding to the type of the message
// supported types are
// - text
// - photo
// - audio
// - video
func (s *State) RegisterTypeRule(tp string, t Transistor) {
	s.TypeRules[tp] = t
}

// Registers a fallback transistor
func (s *State) RegisterFallback(t Transistor) {
	s.FallbackRule = t
}

// This function tests the rules defined for the state and
// returns the next state
//
// The priority order is as follows:
// 1. check cmd rules
// 2. check type rules
// 3. fallback transistor
// On each step if match is found stop and return the state
func (s *State) Test(msg crossbot.Message) (string, error) {
	if msg.Chat.IsCommand() {
		cmd := msg.Chat.Command()
		if nxt, ok := s.CmdRules[cmd]; ok {
			next, err := nxt(msg)
			if err != nil {
				return "", err
			}
			return next, nil
		}
	} else {
		tp := msg.Chat.Type
		if nxt, ok := s.TypeRules[tp]; ok {
			next, err := nxt(msg)
			if err != nil {
				return "", err
			}
			return next, nil
		}
	}

	if s.FallbackRule == nil {
		return "", ErrNoRuleFound
	} else {
		next, err := s.FallbackRule(msg)
		if err != nil {
			return "", err
		}
		return next, nil
	}
}

// This function provides the next state coreesponding to the message
// received by the current state
type Transistor func(msg crossbot.Message) (string, error)

// This is the function that performs action
// on entering or leaving a particular state
type Action func(msg crossbot.Message) error

// Create a new empty state
func MakeState(name string) *State {
	return &State{
		Name:         name,
		Enter:        nil,
		Leave:        nil,
		IsMoved:      False,
		Chain:        "",
		CmdRules:     make(map[string]Transistor),
		TypeRules:    make(map[string]Transistor),
		FallbackRule: nil,
	}
}
