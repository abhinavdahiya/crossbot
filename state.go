package crossbot

import (
	"errors"

	"github.com/abhinavdahiya/crossbot/connector"
)

const (
	ErrNoRuleFound = errors.New("No Rule found to transition the state")
)

// This struct defines any state for the bot
// corresponding to any user
// This also stores transitors for all possible
// moves from this state
//
// make sure the state Name is unique
type State struct {
	Name         string
	Enter, Leave Action
	IsMoved      bool
	Chain        string
	CmdRules     map[string]string
	TypeRules    map[string]string
	FallbackRule string
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

// Registers a command transitor for this state
func (s *State) RegisterCmdRule(cmd string, next string) {
	s.CmdRules[cmd] = next
}

// Registers a transitor corresponding to the type of the message
// supported types are
// - text
// - photo
// - audio
// - video
func (s *State) RegisterTypeRule(tp string, next string) {
	s.TypeRules[tp] = next
}

// Registers a fallback transitor
func (s *State) RegisterFallback(next string) {
	s.FallbackRule = next
}

// This function tests the rules defined for the state and
// returns the next state
//
// The priority order is as follows:
// 1. check cmd rules
// 2. check type rules
// 3. fallback transitor
// On each step if match is found stop and return the state
func (s *State) Test(msg connector.Message) (string, error) {
	if msg.Chat.IsCommand() {
		cmd := msg.Chat.Command()
		if nxt, ok := s.CmdRules[cmd]; ok {
			return nxt, nil
		}
	} else {
		tp := msg.Chat.Type
		if nxt, ok := s.TypeRules[tp]; ok {
			return nxt, nil
		}
	}

	if s.FallbackRule == "" {
		return "", ErrNoRuleFound
	}

	nxt := s.FallbackRule
	return nxt, nil

}

// This is the function that performs action
// on entering or leaving a particular state
type Action func(msg connector.Message) error

// Create a new empty state
func MakeState(name string) *State {
	return &State{
		Name:         name,
		Enter:        nil,
		Leave:        nil,
		IsMoved:      False,
		Chain:        "",
		CmdRules:     make(map[string]string),
		TypeRules:    make(map[string]string),
		FallbackRule: "",
	}
}
