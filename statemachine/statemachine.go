package statemachine

import "fmt"

// State represents the state of the conversation.
type State string

// All possible states for the conversation.
const (
	StateRoot               State = "STATE_ROOT"
	TechOpsFlow             State = "TECH_OPS_FLOW"
	TechOpsState1           State = "TECHOPS_STATE_1"
	TechOpsState2           State = "TECHOPS_STATE_2"
	TechOpsState3           State = "TECHOPS_STATE_3"
	TechOpsState4           State = "TECHOPS_STATE_4"
	TechOpsState5           State = "TECHOPS_STATE_5"
	TechOpsState6           State = "TECHOPS_STATE_6"
	TechOpsState7           State = "TECHOPS_STATE_7"
	TechOpsState8           State = "TECHOPS_STATE_8"
	TechOpsState9           State = "TECHOPS_STATE_9"
	TechOpsState10          State = "TECHOPS_STATE_10"
	TechOpsState11          State = "TECHOPS_STATE_11"
	TechOpsState12          State = "TECHOPS_STATE_12"
	TechOpsState13          State = "TECHOPS_STATE_13"
	TechOpsState14          State = "TECHOPS_STATE_14"
	TechOpsState15          State = "TECHOPS_STATE_15"
	ClientAreaFlow          State = "CLIENT_AREA_FLOW"
	ClientState1            State = "CLIENT_STATE_1"
	ClientState2            State = "CLIENT_STATE_2"
	ClientState3            State = "CLIENT_STATE_3"
	ClientState4            State = "CLIENT_STATE_4"
	ClientState5            State = "CLIENT_STATE_5"
	ClientState6            State = "CLIENT_STATE_6"
	ClientState7            State = "CLIENT_STATE_7"
	ClientState8            State = "CLIENT_STATE_8"
	ClientState9            State = "CLIENT_STATE_9"
	ClientState10           State = "CLIENT_STATE_10"
	ClientState11           State = "CLIENT_STATE_11"
	ClientState12           State = "CLIENT_STATE_12"
	ClientState13           State = "CLIENT_STATE_13"
	ClientState14           State = "CLIENT_STATE_14"
	ClientState15           State = "CLIENT_STATE_15"
	ClientState16           State = "CLIENT_STATE_16"
	ClientState17           State = "CLIENT_STATE_17"
	ClientState18           State = "CLIENT_STATE_18"
)

// Handler is a function that handles a specific state.
type Handler func(userID string, message string) (string, State)

// StateMachine is the state machine for the bot.
type StateMachine struct {
	Handlers map[State]Handler
}

// NewStateMachine creates a new StateMachine.
func NewStateMachine() *StateMachine {
	return &StateMachine{
		Handlers: make(map[State]Handler),
	}
}

// RegisterHandler registers a handler for a specific state.
func (sm *StateMachine) RegisterHandler(state State, handler Handler) {
	sm.Handlers[state] = handler
}

// Handle handles a message for a specific user.
func (sm *StateMachine) Handle(currentState State, userID string, message string) (string, State) {
	handler, ok := sm.Handlers[currentState]
	if !ok {
		return "Desculpe, não entendi. Por favor, tente novamente.", StateRoot
	}
	return handler(userID, message)
}
