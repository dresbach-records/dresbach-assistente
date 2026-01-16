
package whatsapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dresbach/dresbach-assistente/pkg/state"
)

// WebhookHandler processa as requisições do webhook do WhatsApp.
type WebhookHandler struct {
	WhatsAppClient *Client
	StateManager   *state.StateManager
}

// ServeHTTP lida com as requisições HTTP para o webhook.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	var webhookReq map[string]interface{}
	if err := json.Unmarshal(body, &webhookReq); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	messageText := extractMessageText(webhookReq)
	senderPhone := extractSenderPhone(webhookReq)

	if senderPhone == "" || messageText == "" {
		log.Println("Recebida requisição sem remetente ou texto da mensagem.")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Processa a mensagem usando o StateManager, que agora detém toda a lógica de negócios.
	responseText, err := h.StateManager.ProcessMessage(senderPhone, messageText)
	if err != nil {
		log.Printf("Erro ao processar a mensagem pelo StateManager: %v", err)
		responseText = "Desculpe, ocorreu um erro ao processar sua solicitação. Tente novamente."
	}

	// Envia a resposta (gerada pelo StateManager) de volta ao usuário.
	if err := h.WhatsAppClient.SendMessage(senderPhone, responseText); err != nil {
		log.Printf("Erro ao enviar resposta do WhatsApp: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

// As funções auxiliares (extractMessageText, extractSenderPhone) permanecem as mesmas.

func extractMessageText(req map[string]interface{}) string {
	if entry, ok := req["entry"].([]interface{}); ok && len(entry) > 0 {
		if changes, ok := entry[0].(map[string]interface{})["changes"].([]interface{}); ok && len(changes) > 0 {
			if value, ok := changes[0].(map[string]interface{})["value"].(map[string]interface{}); ok {
				if messages, ok := value["messages"].([]interface{}); ok && len(messages) > 0 {
					if text, ok := messages[0].(map[string]interface{})["text"].(map[string]interface{}); ok {
						return text["body"].(string)
					}
				}
			}
		}
	}
	return ""
}

func extractSenderPhone(req map[string]interface{}) string {
	if entry, ok := req["entry"].([]interface{}); ok && len(entry) > 0 {
		if changes, ok := entry[0].(map[string]interface{})["changes"].([]interface{}); ok && len(changes) > 0 {
			if value, ok := changes[0].(map[string]interface{})["value"].(map[string]interface{}); ok {
				if messages, ok := value["messages"].([]interface{}); ok && len(messages) > 0 {
					if from, ok := messages[0].(map[string]interface{})["from"].(string); ok {
						return from
					}
				}
			}
		}
	}
	return ""
}
