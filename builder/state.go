package builder

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

}

// Stores context for a state
func (s *State) SetData(d interface{}) {

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
func (s *State) Test() string {
}

//TODO: msg is Message Struct from crossbot
type Transistor func(msg Message) string
