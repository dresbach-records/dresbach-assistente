
package stripe

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

// AccountProvisioner define a interface para provisionar uma conta após o pagamento.
// Esta é a única dependência que o pacote Stripe tem com o mundo exterior.
type AccountProvisioner interface {
	ProvisionAccount(userID, domain, contactEmail string) error
}

// WebhookHandler lida com os webhooks da Stripe.
type WebhookHandler struct {
	Provisioner         AccountProvisioner
	StripeWebhookSecret string
}

// ServeHTTP processa os eventos de webhook da Stripe.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	// Corrigido: Ler o corpo da requisição para um []byte
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

		// A partir da versão v72 da biblioteca stripe-go, o PaymentIntent está dentro da sessão
		// e precisamos expandi-lo na chamada de criação da sessão. Se não estiver lá,
		// pode ser necessário buscar o PaymentIntent separadamente.
		// Por enquanto, vamos assumir que ele não está diretamente disponível e buscar pelos metadados.
		metadata := session.Metadata
		if metadata == nil {
			// Fallback para o PaymentIntent se os metadados da sessão não estiverem definidos
			// (isso depende de como a sessão de checkout é criada)
			log.Printf("AVISO: Metadados não encontrados diretamente na sessão. Verificando PaymentIntent...")
			// Esta parte pode precisar de mais lógica para buscar o PaymentIntent se não for expandido.
			// Por simplicidade, vamos parar aqui se não encontrarmos os metadados na sessão.
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, ok := metadata["user_id"]
		if !ok {
			log.Printf("ERRO: 'user_id' não encontrado nos metadados da sessão.")
			w.WriteHeader(http.StatusOK) // Responde OK para não ser reenviado
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
	}

	w.WriteHeader(http.StatusOK)
}
