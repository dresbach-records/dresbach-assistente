package products

// Product representa um produto de hospedagem oferecido.
type Product struct {
	ID          string
	Name        string
	Description string
	Type        string
}

// Lista de todos os produtos de hospedagem disponíveis.
var AllProducts = []Product{
	{
		ID:          "prod_TlHnyj6Y4xiHMp",
		Name:        "Hospedagem III",
		Description: "Ideal para projetos robustos até 5 sites com os melhores recursos de velocidade e segurança.",
		Type:        "service",
	},
	{
		ID:          "prod_TlHnHEozuwnuQZ",
		Name:        "Hospedagem II",
		Description: "Ideal para você colocar até 3 sites no ar com recursos avançados de velocidade e segurança.",
		Type:        "service",
	},
	{
		ID:          "prod_TlHmZ8KTx0pe0y",
		Name:        "Hospedagem I",
		Description: "Ideal para colocar o seu site no ar e iniciar o seu projeto em um ambiente de qualidade.",
		Type:        "service",
	},
	{
		ID:          "prod_Tl2cU7tNegEibw",
		Name:        "Hospedagem word press I",
		Description: "Hospedagem I (até 1 site)Ideal para colocar seu site no ar e iniciar seu projeto com qualidade. Inclui 1 site, WordPress gerenciado, migração grátis, criador de sites, 10 GB de armazenamento, 5 emails, largura de banda ilimitada, SSL grátis, backup diário e suporte 24/7.",
		Type:        "service",
	},
    {
        ID:          "prod_TnrOxqyGuMF4EU",
        Name:        "Diagnóstico Técnico + Jurídico",
        Description: "Entrega: até 3 dias úteis",
        Type:        "service",
    },
}
