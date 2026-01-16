
package stripe

import (
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

// Client é um cliente para a API da Stripe.
type Client struct {
	SecretKey string
}

// NewClient cria um novo cliente Stripe.
func NewClient(secretKey string) *Client {
	stripe.Key = secretKey
	return &Client{SecretKey: secretKey}
}

// CreateCheckoutSession cria uma sessão de checkout na Stripe e retorna a URL.
func (c *Client) CreateCheckoutSession(userID, domain string) (string, error) {
	priceID := "price_1P8cysRxx94s4qFjJbAQS5nJ" // SUBSTITUA PELO SEU PRICE ID REAL
	successURL := "https://dresbachhosting.com.br/sucesso" // SUBSTITUA PELA SUA URL DE SUCESSO
	cancelURL := "https://dresbachhosting.com.br/cancelamento"  // SUBSTITUA PELA SUA URL DE CANCELAMENTO

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"boleto",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		// Anexa os metadados ao PaymentIntent, que é o objeto de pagamento subjacente.
		// Isso resolve o erro de compilação "unknown field Metadata".
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"user_id": userID,
				"domain":  domain,
			},
		},
	}

	// Expande o objeto PaymentIntent no retorno da sessão para que o webhook
	// não precise fazer uma chamada extra à API para obter os metadados.
	params.AddExpand("payment_intent")

	s, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("falha ao criar a sessão de checkout na Stripe: %w", err)
	}

	return s.URL, nil
}
