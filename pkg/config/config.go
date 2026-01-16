
package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config armazena todas as configurações da aplicação, carregadas do ambiente.
type Config struct {
	// WhatsApp
	WhatsAppToken           string `envconfig:"WHATSAPP_TOKEN" required:"true"`
	WhatsAppBusinessAccID   string `envconfig:"WHATSAPP_BUSINESS_ACC_ID" required:"true"`
	WhatsAppPhoneNumberID   string `envconfig:"WHATSAPP_PHONE_NUMBER_ID" required:"true"`

	// MongoDB
	MongoURI string `envconfig:"MONGO_URI" required:"true"`

	// WHM
	WHMHost     string `envconfig:"WHM_HOST" required:"true"`
	WHMAPIToken string `envconfig:"WHM_API_TOKEN" required:"true"`

	// Stripe
	StripeKey           string `envconfig:"STRIPE_KEY" required:"true"`
	StripeWebhookSecret string `envconfig:"STRIPE_WEBHOOK_SECRET" required:"true"`
}

// New carrega a configuração a partir de variáveis de ambiente.
func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("falha ao carregar configuração do ambiente: %w", err)
	}
	return &cfg, nil
}
