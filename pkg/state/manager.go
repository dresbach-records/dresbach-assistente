
package state

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dresbach/dresbach-assistente/pkg/database"
	"github.com/dresbach/dresbach-assistente/pkg/stripe"
	"github.com/dresbach/dresbach-assistente/pkg/whm"
)

// State representa o estado atual de uma conversa de usuário.
type State string

// Definição dos estados da conversa
const (
	StateInitial        State = "INITIAL"
	StateAwaitingOption State = "AWAITING_OPTION"

	// Fluxo: Área do Cliente
	StateClientLogin State = "CLIENT_LOGIN"

	// Fluxo: Suporte
	StateSupport State = "SUPPORT"

	// Fluxo: Tech Ops (Contratar um serviço)
	StateTechOpsStart                State = "TECHOPS_START"
	StateTechOpsLgpdConfirm          State = "TECHOPS_LGPD_CONFIRM"
	StateTechOpsPainPoint            State = "TECHOPS_PAIN_POINT"
	StateTechOpsHistory              State = "TECHOPS_HISTORY"
	StateTechOpsPreAnalysisStart     State = "TECHOPS_PRE_ANALYSIS_START"
	StateTechOpsPreAnalysisRepo      State = "TECHOPS_PRE_ANALYSIS_REPO"
	StateTechOpsPreAnalysisRepoLink  State = "TECHOPS_PRE_ANALYSIS_REPO_LINK"
	StateTechOpsPreAnalysisSite      State = "TECHOPS_PRE_ANALYSIS_SITE"
	StateTechOpsPreAnalysisSiteLink  State = "TECHOPS_PRE_ANALYSIS_SITE_LINK"
	StateTechOpsPreAnalysisProblem   State = "TECHOPS_PRE_ANALYSIS_PROBLEM"
	StateTechOpsClassification       State = "TECHOPS_CLASSIFICATION"

	// Novos estados para o fluxo de domínio
	StateTechOpsAskDomain        State = "TECHOPS_ASK_DOMAIN"
	StateTechOpsCheckDomainOwner State = "TECHOPS_CHECK_DOMAIN_OWNER"
	StateTechOpsHandleTransfer   State = "TECHOPS_HANDLE_TRANSFER"
	StateTechOpsHandleRegister   State = "TECHOPS_HANDLE_REGISTER"

	StateAwaitingPayment State = "AWAITING_PAYMENT"

	// Estado final (exemplo)
	StateFinalizing State = "FINALIZING"
)


// StateManager gerencia a lógica de estados da conversa, agora com persistência.
type StateManager struct {
	dbStore      *database.MongoStore
	whmClient    *whm.Client
	stripeClient *stripe.Client // Adicionado o cliente Stripe
}

// NewManager cria um novo StateManager com todas as dependências necessárias.
func NewManager(dbStore *database.MongoStore, whmClient *whm.Client, stripeClient *stripe.Client) *StateManager {
	return &StateManager{
		dbStore:      dbStore,
		whmClient:    whmClient,
		stripeClient: stripeClient, // Injetando o cliente Stripe
	}
}

// ProcessMessage processa a mensagem e conduz a conversa com base no estado.
func (sm *StateManager) ProcessMessage(userID, messageText string) (string, error) {
	ctx := context.Background()

	// Carrega a sessão do banco de dados.
	dbSession, err := sm.dbStore.LoadSession(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("falha ao carregar sessão no StateManager: %w", err)
	}
	session := copyFromDBSession(dbSession)

	var response string
	normalizedInput := strings.ToLower(strings.TrimSpace(messageText))

	switch session.State {
	// ... (casos anteriores permanecem os mesmos)

	case StateTechOpsHandleTransfer, StateTechOpsHandleRegister:
		if normalizedInput == "ok" {
			// Gera o link de checkout da Stripe
			checkoutURL, err := sm.stripeClient.CreateCheckoutSession(session.UserID, session.Domain)
			if err != nil {
				log.Printf("ERRO: Falha ao criar a sessão de checkout da Stripe: %v", err)
				response = "Desculpe, não consegui gerar o link de pagamento. Por favor, tente novamente mais tarde."
			} else {
				response = fmt.Sprintf("Tudo pronto! Para finalizar, efetue o pagamento através deste link seguro: %s", checkoutURL)
				session.State = StateAwaitingPayment // Muda o estado para aguardar a confirmação do pagamento
			}
		} else {
			response = "Por favor, digite 'OK' para confirmar e continuar."
		}

	// ... (outros casos)
	default:
		// Resetar para o estado inicial se algo der errado ou um estado não for tratado.
		response = fmt.Sprintf("Ocorreu um erro (Estado não implementado: %s). Reiniciando a conversa.", session.State)
		session.State = StateInitial

	}

	// Salva a sessão atualizada de volta ao banco de dados.
	if err := sm.dbStore.SaveSession(ctx, copyToDBSession(session)); err != nil {
		return "", fmt.Errorf("falha ao salvar sessão no StateManager: %w", err)
	}

	log.Printf("Sessão salva para o usuário %s no DB: Estado %s", session.UserID, session.State)

	return response, nil
}

// ProvisionAccount é a função que será chamada pelo webhook da Stripe
func (sm *StateManager) ProvisionAccount(userID, domain, contactEmail string) error {
	// 1. Gera uma senha aleatória para a nova conta
	// (Em um cenário real, use uma biblioteca para gerar senhas seguras)
	password := "senhaSuperSegura123!"

	// 2. Define os parâmetros para a criação da conta no WHM
	params := whm.CreateAccountParams{
		Username:     generateUsername(domain), // Gera um nome de usuário a partir do domínio
		Domain:       domain,
		Plan:         "Dresbach-Start", // SUBSTITUA pelo seu nome de plano real
		Password:     password,
		ContactEmail: contactEmail,
	}

	// 3. Chama a função do cliente WHM para criar a conta
	_, err := sm.whmClient.CreateAccount(params)
	if err != nil {
		log.Printf("ERRO AO PROVISIONAR CONTA: Falha ao criar conta no WHM para o usuário %s: %v", userID, err)
		return fmt.Errorf("falha ao criar conta no WHM: %w", err)
	}

	log.Printf("SUCESSO: Conta para o domínio %s (usuário %s) provisionada com sucesso!", domain, userID)

	// 4. (Próximo passo) Enviar a confirmação e os dados de acesso para o usuário via WhatsApp
	return nil
}

// generateUsername cria um nome de usuário de 8 caracteres a partir do domínio para evitar colisões.
func generateUsername(domain string) string {
	// Remove tld, www, etc.
	username := strings.Split(domain, ".")[0]
	if len(username) > 8 {
		return username[:8]
	}
	return username
}


// --- Funções auxiliares e Structs ---

// UserSession é a representação da sessão usada internamente pelo StateManager.
type UserSession struct {
	UserID      string
	State       State
	Domain      string // Campo para armazenar o domínio
	PreAnalysis PreAnalysisData
}

// PreAnalysisData é a representação interna dos dados de pré-análise.
type PreAnalysisData struct {
	RepoURL            string
	SystemURL          string
	ProblemDescription string
}

// copyToDBSession converte a UserSession do state para a Session do database.
func copyToDBSession(session *UserSession) *database.Session {
	return &database.Session{
		UserID: session.UserID,
		State:  string(session.State),
		Domain: session.Domain,
		PreAnalysis: database.PreAnalysisData{
			RepoURL:            session.PreAnalysis.RepoURL,
			SystemURL:          session.PreAnalysis.SystemURL,
			ProblemDescription: session.PreAnalysis.ProblemDescription,
		},
	}
}

// copyFromDBSession converte a Session do database para a UserSession do state.
func copyFromDBSession(dbSession *database.Session) *UserSession {
	return &UserSession{
		UserID: dbSession.UserID,
		State:  State(dbSession.State),
		Domain: dbSession.Domain,
		PreAnalysis: PreAnalysisData{
			RepoURL:            dbSession.PreAnalysis.RepoURL,
			SystemURL:          dbSession.PreAnalysis.SystemURL,
			ProblemDescription: dbSession.PreAnalysis.ProblemDescription,
		},
	}
}
