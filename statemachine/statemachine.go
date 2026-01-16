package statemachine

// State represents the state of the conversation.
type State string

// All possible states for the conversation.
const (
	StateRoot                     State = "STATE_ROOT"
	TechOpsFlow                   State = "TECH_OPS_FLOW"
	TechOpsQualifyLeadHasSystem   State = "TECHOPS_QUALIFY_LEAD_HAS_SYSTEM"
	TechOpsQualifyLeadCollectsData State = "TECHOPS_QUALIFY_LEAD_COLLECTS_DATA"
	TechOpsPrincipalProblem       State = "TECHOPS_PRINCIPAL_PROBLEM"
	TechOpsPaidDiagnosis          State = "TECHOPS_PAID_DIAGNOSIS"
	TechOpsPayment                State = "TECHOPS_PAYMENT"
	TechOpsPaymentValidated       State = "TECHOPS_PAYMENT_VALIDATED"
	TechOpsTechnicalDiagnosis     State = "TECHOPS_TECHNICAL_DIAGNOSIS"
	TechOpsAutomaticFeedback      State = "TECHOPS_AUTOMATIC_FEEDBACK"
	TechOpsSpecializedConsulting  State = "TECHOPS_SPECIALIZED_CONSULTING"

	ClientAreaFlow      State = "CLIENT_AREA_FLOW"
	ClientLogin         State = "CLIENT_LOGIN"
	ClientMainMenu      State = "CLIENT_MAIN_MENU"
	ClientServices      State = "CLIENT_SERVICES"
	ClientServicesDetails State = "CLIENT_SERVICES_DETAILS"
	ClientDomains       State = "CLIENT_DOMAINS"
	ClientInvoices      State = "CLIENT_INVOICES"
	ClientTickets       State = "CLIENT_TICKETS"
	ClientViewTickets   State = "CLIENT_VIEW_TICKETS"
	ClientOpenTicket    State = "CLIENT_OPEN_TICKET"
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
