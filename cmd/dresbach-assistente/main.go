
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dresbach/dresbach-assistente/pkg/config"
	"github.com/dresbach/dresbach-assistente/pkg/database"
	"github.com/dresbach/dresbach-assistente/pkg/state"
	"github.com/dresbach/dresbach-assistente/pkg/stripe"
	"github.com/dresbach/dresbach-assistente/pkg/whatsapp"
	"github.com/dresbach/dresbach-assistente/pkg/whm"

	"github.com/joho/godotenv" // Importa a biblioteca
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível encontrar o arquivo .env, as variáveis de ambiente devem ser setadas manualmente.")
	}

	ctx := context.Background()

	// 1. Carrega as configurações de ambiente
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Erro ao carregar a configuração: %v", err)
	}

	// 2. Inicializa a conexão com o MongoDB
	dbStore, err := database.NewMongoStore(ctx, cfg.MongoURI, "dresbach_assistant", "sessions")
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer dbStore.Close(ctx)

	// 3. Inicializa os clientes dos serviços externos
	whmClient := whm.NewClient(cfg.WHMHost, cfg.WHMAPIToken)
	stripeClient := stripe.NewClient(cfg.StripeKey)
	whatsappClient := whatsapp.NewClient(cfg.WhatsAppToken, cfg.WhatsAppBusinessAccID, cfg.WhatsAppPhoneNumberID)

	// 4. Inicializa o StateManager, injetando todas as dependências
	stateManager := state.NewManager(dbStore, whmClient, stripeClient, whatsappClient)

	log.Println("Iniciando o servidor Dresbach Assistente na porta 8080...")

	// 5. Cria os Handlers HTTP
	whatsappWebhookHandler := &whatsapp.WebhookHandler{
		WhatsAppClient: whatsappClient,
		StateManager:   stateManager,
	}

	// O WebhookHandler do Stripe agora só precisa do Provisioner
	stripeWebhookHandler := &stripe.WebhookHandler{
		Provisioner:         stateManager, // stateManager implementa a interface AccountProvisioner
		StripeWebhookSecret: cfg.StripeWebhookSecret,
	}

	// 6. Registra os handlers e inicia o servidor
	http.Handle("/whatsapp-webhook", whatsappWebhookHandler)
	http.Handle("/stripe-webhook", stripeWebhookHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
