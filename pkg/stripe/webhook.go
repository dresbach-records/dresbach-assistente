
package stripe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dresbach/dresbach-assistente/pkg/whatsapp"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

// AccountProvisioner define a interface para provisionar uma conta após o pagamento.
type AccountProvisioner interface {
	ProvisionAccount(userID, domain, contactEmail string) error
}

// WebhookHandler lida com os webhooks da Stripe.
type WebhookHandler struct {
	Provisioner         AccountProvisioner
	WhatsAppClient      *whatsapp.Client
	StripeWebhookSecret string
}

// ServeHTTP processa os eventos de webhook da Stripe.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERRO: Falha ao ler o corpo do webhook da Stripe: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.StripeWebhookSecret)
	if err != nil {
		log.Printf("ERRO: Falha na verificação da assinatura do webhook da Stripe: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("ERRO: Falha ao decodificar a sessão do evento da Stripe: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verifica se o objeto PaymentIntent foi expandido e incluído.
		if session.PaymentIntent == nil {
			log.Printf("ERRO CRÍTICO: O objeto PaymentIntent não foi encontrado no webhook da Stripe. Certifique-se de expandi-lo com 'payment_intent' na criação da sessão.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Extrai os metadados do PaymentIntent, não da sessão.
		metadata := session.PaymentIntent.Metadata
		userID, ok := metadata["user_id"]
		if !ok {
			log.Printf("ERRO: 'user_id' não encontrado nos metadados do PaymentIntent do webhook da Stripe.")
			w.WriteHeader(http.StatusOK)
			return
		}
		domain, _ := metadata["domain"]
		contactEmail := session.CustomerEmail

		log.Printf("PAGAMENTO BEM-SUCEDIDO recebido para o usuário: %s, domínio: %s", userID, domain)

		if err := h.Provisioner.ProvisionAccount(userID, domain, contactEmail); err != nil {
			log.Printf("ERRO CRÍTICO: O pagamento foi recebido, mas o provisionamento FALHOU para o usuário %s: %v", userID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("Ótima notícia! ✅\n\nSeu pagamento foi confirmado e sua conta para o domínio `%s` foi criada com sucesso!\n\nEm breve você receberá seus dados de acesso.", domain)
		if err := h.WhatsAppClient.SendMessage(userID, msg); err != nil {
			log.Printf("ERRO: Falha ao enviar mensagem de confirmação de provisionamento para o usuário %s: %v", userID, err)
		}
	}

	w.WriteHeader(http.StatusOK)
}
