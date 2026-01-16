package whatsapp

import (
	"bytes"
	"fmt"
	"net/http"
)

// Client é um cliente para interagir com a API do WhatsApp
type Client struct {
	Token      string
	BusinessID string
	PhoneNumberID string
	HTTPClient *http.Client
}

// NewClient cria um novo cliente do WhatsApp
func NewClient(token, businessID, phoneNumberID string) *Client {
	return &Client{
		Token:      token,
		BusinessID: businessID,
		PhoneNumberID: phoneNumberID,
		HTTPClient: &http.Client{},
	}
}

// SendMessage envia uma mensagem de texto para um destinatário
func (c *Client) SendMessage(to, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", c.PhoneNumberID)

	payload := []byte(fmt.Sprintf(`{ 
		"messaging_product": "whatsapp", 
		"to": "%s", 
		"type": "text", 
		"text": { "body": "%s" } 
	}`, to, message))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		// Lidar com respostas de erro da API do WhatsApp
		return fmt.Errorf("falha ao enviar mensagem do WhatsApp, status: %s", resp.Status)
	}

	return nil
}
