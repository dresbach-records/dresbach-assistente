
package stripe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dresbach/dresbach-assistente/pkg/state"
	"github.com/dresbach/dresbach-assistente/pkg/whatsapp"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

// WebhookHandler lida com os webhooks da Stripe.
type WebhookHandler struct {
	StateManager      *state.StateManager
	WhatsAppClient    *whatsapp.Client
	StripeWebhookSecret string
}

// HandleWebhook processa os eventos de webhook da Stripe.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERRO: Falha ao ler o corpo do webhook da Stripe: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Valida a assinatura do webhook para garantir que veio da Stripe
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.StripeWebhookSecret)
	if err != nil {
		log.Printf("ERRO: Falha na verificação da assinatura do webhook da Stripe: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Processa apenas o evento de checkout bem-sucedido
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("ERRO: Falha ao decodificar a sessão do evento da Stripe: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Extrai os metadados que associamos
		userID, ok := session.Metadata["user_id"]
		if !ok {
			log.Printf("ERRO: 'user_id' não encontrado nos metadados do webhook da Stripe.")
			// Retorna 200 para a Stripe não reenviar, pois não há como corrigir.
			w.WriteHeader(http.StatusOK)
			return
		}
		domain, _ := session.Metadata["domain"]
		contactEmail := session.CustomerEmail // Captura o email do cliente

		log.Printf("PAGAMENTO BEM-SUCEDIDO recebido para o usuário: %s, domínio: %s", userID, domain)

		// Chama a função de provisionamento
		if err := h.StateManager.ProvisionAccount(userID, domain, contactEmail); err != nil {
			log.Printf("ERRO CRÍTICO: O pagamento foi recebido, mas o provisionamento FALHOU para o usuário %s: %v", userID, err)
			// (Opcional) Enviar uma notificação para a equipe interna sobre a falha.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Envia a confirmação para o cliente via WhatsApp
		msg := fmt.Sprintf("Ótima notícia! ✅\n\nSeu pagamento foi confirmado e sua conta para o domínio `%s` foi criada com sucesso!\n\nEm breve você receberá seus dados de acesso.\n(Lógica de envio de senha a ser implementada)", domain)
		if err := h.WhatsAppClient.SendMessage(userID, msg); err != nil {
			log.Printf("ERRO: Falha ao enviar mensagem de confirmação de provisionamento para o usuário %s: %v", userID, err)
		}
	}

	w.WriteHeader(http.StatusOK)
}
