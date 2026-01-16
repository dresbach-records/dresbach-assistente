
package whm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client é um cliente para a API do WHM.
type Client struct {
	Host       string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient cria um novo cliente WHM.
func NewClient(host, apiToken string) *Client {
	return &Client{
		Host:       host,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}

// CreateAccountParams define os parâmetros para a criação de uma conta no WHM.
type CreateAccountParams struct {
	Username string
	Domain   string
	Plan     string
	Password string
	ContactEmail string
}

// WHMResponseMetadata contém metadados comuns da resposta da API.
type WHMResponseMetadata struct {
	Command string `json:"command"`
	Reason  string `json:"reason"`
	Result  int    `json:"result"`
}

// CreateAccountResponse define a estrutura da resposta bem-sucedida de 'createacct'.
type CreateAccountResponse struct {
	Metadata WHMResponseMetadata `json:"metadata"`
}


// CreateAccount chama a API 'createacct' do WHM para provisionar uma nova conta.
func (c *Client) CreateAccount(params CreateAccountParams) (*CreateAccountResponse, error) {
	// 1. Monta a URL da API com os parâmetros necessários.
	query := url.Values{}
	query.Set("api.version", "1")
	query.Set("username", params.Username)
	query.Set("domain", params.Domain)
	query.Set("plan", params.Plan)
	query.Set("password", params.Password)
	query.Set("contactemail", params.ContactEmail)

	apiURL := fmt.Sprintf("https://%s:2087/json-api/createacct?%s", c.Host, query.Encode())

	// 2. Cria a requisição HTTP GET.
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição para o WHM: %w", err)
	}

	// 3. Adiciona o header de autorização.
	req.Header.Add("Authorization", fmt.Sprintf("whm %s", c.APIToken))

	// 4. Executa a requisição.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar requisição para o WHM: %w", err)
	}
	defer resp.Body.Close()

	// 5. Lê e decodifica a resposta.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta do WHM: %w", err)
	}

	var whmResponse CreateAccountResponse
	if err := json.Unmarshal(body, &whmResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON do WHM: %w (resposta: %s)", err, string(body))
	}

	// 6. Verifica se a API retornou um erro lógico.
	if whmResponse.Metadata.Result == 0 {
		return nil, fmt.Errorf("erro da API WHM: %s", whmResponse.Metadata.Reason)
	}

	// 7. Retorna a resposta bem-sucedida.
	return &whmResponse, nil
}
