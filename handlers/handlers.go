package handlers

import (
	"fmt"

	"github.com/dresbach/project/statemachine"
)

// RegisterHandlers registers all the handlers for the state machine.
func RegisterHandlers(sm *statemachine.StateMachine) {
	sm.RegisterHandler(statemachine.StateRoot, handleRoot)

	// Tech Ops Flow
	sm.RegisterHandler(statemachine.TechOpsFlow, handleTechOpsFlow)
	sm.RegisterHandler(statemachine.TechOpsState1, handleTechOpsState1)
	sm.RegisterHandler(statemachine.TechOpsState2, handleTechOpsState2)
	sm.RegisterHandler(statemachine.TechOpsState3, handleTechOpsState3)
	sm.RegisterHandler(statemachine.TechOpsState4, handleTechOpsState4)
	sm.RegisterHandler(statemachine.TechOpsState5, handleTechOpsState5)
	sm.RegisterHandler(statemachine.TechOpsState6, handleTechOpsState6)
	sm.RegisterHandler(statemachine.TechOpsState7, handleTechOpsState7)
	sm.RegisterHandler(statemachine.TechOpsState8, handleTechOpsState8)
	sm.RegisterHandler(statemachine.TechOpsState9, handleTechOpsState9)
	sm.RegisterHandler(statemachine.TechOpsState10, handleTechOpsState10)
	sm.RegisterHandler(statemachine.TechOpsState11, handleTechOpsState11)
	sm.RegisterHandler(statemachine.TechOpsState12, handleTechOpsState12)
	sm.RegisterHandler(statemachine.TechOpsState13, handleTechOpsState13)
	sm.RegisterHandler(statemachine.TechOpsState14, handleTechOpsState14)
	sm.RegisterHandler(statemachine.TechOpsState15, handleTechOpsState15)

	// Client Area Flow
	sm.RegisterHandler(statemachine.ClientAreaFlow, handleClientAreaFlow)
	sm.RegisterHandler(statemachine.ClientState1, handleClientState1)
	// ... register other client area handlers
}

func handleRoot(userID, message string) (string, statemachine.State) {
	switch message {
	case "1":
		return "Para começarmos, qual é o tipo do seu negócio?", statemachine.TechOpsState1
	case "2":
		return "Para acessar a área do cliente, por favor, informe seu CPF.", statemachine.ClientState1
	default:
		return "👋 Olá! Bem-vindo à *Dresbach Hosting do Brasil*.\n\nEscolha como podemos te ajudar:\n1️⃣ Tech Ops — Consultoria Especializada\n2️⃣ Área do Cliente — Hospedagem\n\nDigite o número da opção desejada.", statemachine.StateRoot
	}
}

// Tech Ops Flow Handlers

func handleTechOpsFlow(userID, message string) (string, statemachine.State) {
	// This state is a routing state, so it should not be reached.
	return "Erro interno.", statemachine.StateRoot
}

func handleTechOpsState1(userID, message string) (string, statemachine.State) {
	// Placeholder for saving business type
	return "Hoje você já possui site ou sistema em funcionamento?\n(Sim / Não)", statemachine.TechOpsState2
}

func handleTechOpsState2(userID, message string) (string, statemachine.State) {
	switch message {
	case "Sim":
		return "Esse sistema está em produção ou em desenvolvimento?", statemachine.TechOpsState3
	case "Não":
		return "Você pretende criar um site institucional ou um sistema com login e dados?", statemachine.TechOpsState3
	default:
		return "Por favor, responda com Sim ou Não.", statemachine.TechOpsState2
	}
}

func handleTechOpsState3(userID, message string) (string, statemachine.State) {
	// Placeholder for saving system status
	return "Esse sistema coleta ou armazena dados de clientes ou usuários?\n(Sim / Não)", statemachine.TechOpsState4
}

func handleTechOpsState4(userID, message string) (string, statemachine.State) {
	// Placeholder for saving data collection status
	return "Qual é hoje sua principal preocupação?\n1️⃣ Segurança / invasão\n2️⃣ LGPD / jurídico\n3️⃣ Estrutura técnica / performance\n4️⃣ Não sei por onde começar", statemachine.TechOpsState5
}

func handleTechOpsState5(userID, message string) (string, statemachine.State) {
	// Placeholder for saving main concern
	return "Você já teve algum problema como invasão, vazamento ou queda?\n(Sim / Não)", statemachine.TechOpsState6
}

func handleTechOpsState6(userID, message string) (string, statemachine.State) {
	// Placeholder for saving problem history
	return "Com base no que você informou, o mais indicado é um diagnóstico técnico e jurídico inicial.\nPosso te explicar como funciona?", statemachine.TechOpsState7
}

func handleTechOpsState7(userID, message string) (string, statemachine.State) {
	return "✔ Avaliamos arquitetura técnica\n✔ Identificamos riscos de segurança\n✔ Analisamos LGPD e contratos\n✔ Indicamos próximos passos\n\nIsso evita retrabalho e custos desnecessários.", statemachine.TechOpsState8
}

func handleTechOpsState8(userID, message string) (string, statemachine.State) {
	return "🔍 Diagnóstico Técnico + Jurídico\n💰 R$ 297,00\n⏱ Entrega: até 3 dias úteis\n\nDeseja seguir?\n(Sim / Não)", statemachine.TechOpsState9
}

func handleTechOpsState9(userID, message string) (string, statemachine.State) {
	switch message {
	case "Sim":
		return "Forma de pagamento:\n1️⃣ PIX\n2️⃣ Cartão\n\nPIX:\nChave: 51981446019\nBeneficiário: Dresbach Hosting do Brasil LTDA\n\nEnvie o comprovante após o pagamento.", statemachine.TechOpsState10
	case "Não":
		// Handle the case where the user does not want to proceed
		return "Entendido. Se precisar de algo mais, é só chamar.", statemachine.StateRoot
	default:
		return "Por favor, responda com Sim ou Não.", statemachine.TechOpsState9
	}
}

func handleTechOpsState10(userID, message string) (string, statemachine.State) {
	// Placeholder for handling payment method selection
	// Placeholder for integrating with payment API
	return "Pagamento confirmado ✅", statemachine.TechOpsState11
}

func handleTechOpsState11(userID, message string) (string, statemachine.State) {
	return "Escolha um horário disponível:\n1️⃣ 10:00\n2️⃣ 15:00\n3️⃣ Próximo dia útil 09:00", statemachine.TechOpsState12
}

func handleTechOpsState12(userID, message string) (string, statemachine.State) {
	// Placeholder for saving the chosen time
	// Placeholder for integrating with scheduling API
	return "Consultoria agendada com sucesso.\nUm especialista entrará em contato.", statemachine.TechOpsState13
}

func handleTechOpsState13(userID, message string) (string, statemachine.State) {
	return "Antes de finalizar, poderia nos dizer:\nComo foi sua experiência até aqui?\n(Ótima / Boa / Pode melhorar)", statemachine.TechOpsState14
}

func handleTechOpsState14(userID, message string) (string, statemachine.State) {
	// Placeholder for saving feedback
	// Transfer to human operator
	// Tag: TECHOPS_DIAGNOSTICO_PAGO_AGENDADO
	return "Obrigado pelo seu feedback! Um de nossos especialistas entrará em contato em breve.", statem-achine.TechOpsState15
}

func handleTechOpsState15(userID, message string) (string, statemachine.State) {
	// This state is for transferring to a human operator.
	// The logic for the transfer should be implemented here.
	return "", statemachine.StateRoot // Or a specific state after human interaction
}

// Client Area Flow Handlers

func handleClientAreaFlow(userID, message string) (string, statemachine.State) {
	// This state is a routing state, so it should not be reached.
	return "Erro interno.", statemachine.StateRoot
}

func handleClientState1(userID, message string) (string, statemachine.State) {
	// Placeholder for CPF validation
	return "Agora, digite sua senha.", statemachine.ClientState2
}

// ... other client area handlers
