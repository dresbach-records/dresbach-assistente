package gemini

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client é um wrapper para o modelo generativo do Gemini.
type Client struct {
	model *genai.GenerativeModel
}

// NewClient cria um novo cliente Gemini.
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()

	// Primeiro, crie o cliente com a chave da API.
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	// Em seguida, obtenha o modelo generativo do cliente.
	model := client.GenerativeModel("gemini-pro")

	return &Client{model: model}, nil
}

// GenerateText gera uma resposta de texto a partir de um prompt.
func (c *Client) GenerateText(prompt string) (string, error) {
	ctx := context.Background()
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var responseText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					responseText += string(txt)
				}
			}
		}
	}

	if responseText == "" {
		log.Println("Resposta do Gemini vazia ou em formato inesperado.")
		return "Desculpe, não consegui processar sua solicitação no momento.", nil
	}

	return responseText, nil
}
