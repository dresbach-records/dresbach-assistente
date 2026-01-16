package handlers

import (
	"github.com/dresbach/project/statemachine"
)

// RegisterHandlers registers all the handlers for the state machine.
func RegisterHandlers(sm *statemachine.StateMachine) {
	sm.RegisterHandler(statemachine.StateRoot, handleRoot)

	// Tech Ops Flow
	sm.RegisterHandler(statemachine.TechOpsFlow, handleTechOpsFlow)
	sm.RegisterHandler(statemachine.TechOpsQualifyLeadHasSystem, handleTechOpsQualifyLeadHasSystem)
	sm.RegisterHandler(statemachine.TechOpsQualifyLeadCollectsData, handleTechOpsQualifyLeadCollectsData)
	sm.RegisterHandler(statemachine.TechOpsPrincipalProblem, handleTechOpsPrincipalProblem)
	sm.RegisterHandler(statemachine.TechOpsPaidDiagnosis, handleTechOpsPaidDiagnosis)
	sm.RegisterHandler(statemachine.TechOpsPayment, handleTechOpsPayment)
	sm.RegisterHandler(statemachine.TechOpsPaymentValidated, handleTechOpsPaymentValidated)
	sm.RegisterHandler(statemachine.TechOpsTechnicalDiagnosis, handleTechOpsTechnicalDiagnosis)
	sm.RegisterHandler(statemachine.TechOpsAutomaticFeedback, handleTechOpsAutomaticFeedback)
	sm.RegisterHandler(statemachine.TechOpsSpecializedConsulting, handleTechOpsSpecializedConsulting)

	// Client Area Flow
	sm.RegisterHandler(statemachine.ClientAreaFlow, handleClientAreaFlow)
	sm.RegisterHandler(statemachine.ClientLogin, handleClientLogin)
	sm.RegisterHandler(statemachine.ClientMainMenu, handleClientMainMenu)
	sm.RegisterHandler(statemachine.ClientServices, handleClientServices)
	sm.RegisterHandler(statemachine.ClientServicesDetails, handleClientServicesDetails)
	sm.RegisterHandler(statemachine.ClientDomains, handleClientDomains)
	sm.RegisterHandler(statemachine.ClientInvoices, handleClientInvoices)
	sm.RegisterHandler(statemachine.ClientTickets, handleClientTickets)
	sm.RegisterHandler(statemachine.ClientViewTickets, handleClientViewTickets)
	sm.RegisterHandler(statemachine.ClientOpenTicket, handleClientOpenTicket)
}

func handleRoot(userID, message string) (string, statemachine.State) {
	switch message {
	case "1":
		return "Você tem sistema ou site?", statemachine.TechOpsQualifyLeadHasSystem
	case "2":
		return "Para acessar a área do cliente, por favor, informe seu login.", statemachine.ClientLogin
	default:
		return "👋 Olá! Bem-vindo à Dresbach Hosting.\n\nEscolha uma opção:\n1️⃣ Tech Ops — Consultoria Especializada\n2️⃣ Área do Cliente — Hospedagem", statemachine.StateRoot
	}
}

// Tech Ops Flow Handlers

func handleTechOpsFlow(userID, message string) (string, statemachine.State) {
	// This state is a routing state, so it should not be reached.
	return "Erro interno.", statemachine.StateRoot
}

func handleTechOpsQualifyLeadHasSystem(userID, message string) (string, statemachine.State) {
	// Placeholder for saving if the user has a system or site
	return "Coleta dados?", statemachine.TechOpsQualifyLeadCollectsData
}

func handleTechOpsQualifyLeadCollectsData(userID, message string) (string, statemachine.State) {
	// Placeholder for saving if the user collects data
	return "Qual o principal problema?", statemachine.TechOpsPrincipalProblem
}

func handleTechOpsPrincipalProblem(userID, message string) (string, statemachine.State) {
	// Placeholder for saving the principal problem
	return "Diagnóstico Pago: R$ 297", statemachine.TechOpsPaidDiagnosis
}

func handleTechOpsPaidDiagnosis(userID, message string) (string, statemachine.State) {
	// Placeholder for handling paid diagnosis
	return "Pagamento", statemachine.TechOpsPayment
}

func handleTechOpsPayment(userID, message string) (string, statemachine.State) {
	// Placeholder for handling payment
	return "Pagamento Validado", statemachine.TechOpsPaymentValidated
}

func handleTechOpsPaymentValidated(userID, message string) (string, statemachine.State) {
	// Placeholder for handling payment validation
	return "Diagnóstico Técnico + Jurídico", statemachine.TechOpsTechnicalDiagnosis
}

func handleTechOpsTechnicalDiagnosis(userID, message string) (string, statemachine.State) {
	// Placeholder for handling technical diagnosis
	return "Feedback Automático", statemachine.TechOpsAutomaticFeedback
}

func handleTechOpsAutomaticFeedback(userID, message string) (string, statemachine.State) {
	// Placeholder for handling automatic feedback
	return "Consultoria Especializada", statemachine.TechOpsSpecializedConsulting
}

func handleTechOpsSpecializedConsulting(userID, message string) (string, statemachine.State) {
	// End of the Tech Ops flow
	return "Obrigado por usar a consultoria especializada.", statemachine.StateRoot
}

// Client Area Flow Handlers

func handleClientAreaFlow(userID, message string) (string, statemachine.State) {
	// This state is a routing state, so it should not be reached.
	return "Erro interno.", statemachine.StateRoot
}

func handleClientLogin(userID, message string) (string, statemachine.State) {
	// Placeholder for handling login
	return "Serviços, Domínios, Faturas", statemachine.ClientMainMenu
}

func handleClientMainMenu(userID, message string) (string, statemachine.State) {
	switch message {
	case "Serviços":
		return "Detalhes", statemachine.ClientServices
	case "Domínios":
		return "Gerenciar Domínios", statemachine.ClientDomains
	case "Faturas":
		return "Ver Faturas", statemachine.ClientInvoices
	default:
		return "Opção inválida. Por favor, escolha entre Serviços, Domínios ou Faturas.", statemachine.ClientMainMenu
	}
}

func handleClientServices(userID, message string) (string, statemachine.State) {
	// Placeholder for handling services
	return "Acessar cPanel, Acessar WHM, Acessar Webmail", statemachine.ClientServicesDetails
}

func handleClientServicesDetails(userID, message string) (string, statemachine.State) {
	// Placeholder for handling service details
	return "Gestão de Hospedagem", statemachine.StateRoot
}

func handleClientDomains(userID, message string) (string, statemachine.State) {
	// Placeholder for handling domains
	return "Ver Titas", statemachine.ClientTickets
}

func handleClientInvoices(userID, message string) (string, statemachine.State) {
	// Placeholder for handling invoices
	return "Ver Faturas", statemachine.ClientTickets
}

func handleClientTickets(userID, message string) (string, statemachine.State) {
	return "Ver Tickets, Abrir Ticket", statemachine.ClientTickets
}

func handleClientViewTickets(userID, message string) (string, statemachine.State) {
	// Placeholder for viewing tickets
	return "Gestão de Hospedagem", statemachine.StateRoot
}

func handleClientOpenTicket(userID, message string) (string, statemachine.State) {
	// Placeholder for opening a ticket
	return "Gestão de Hospedagem", statemachine.StateRoot
}
