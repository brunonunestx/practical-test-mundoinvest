package clients

type Client struct {
	name          string
	email         string
	clientType    string
	status        string
	priority      string
	heritageValue int
}

type ClientResponse struct {
	Name          string  `json:"nome"`
	Email         string  `json:"email"`
	ClientType    string  `json:"tipo_solicitacao"`
	Status        string  `json:"status"`
	Priority      string  `json:"prioridade"`
	HeritageValue float64 `json:"valor_patrimonio"`
}

func (c *Client) toResponse() *ClientResponse {
	return &ClientResponse{
		Name:          c.name,
		Email:         c.email,
		ClientType:    c.clientType,
		Status:        c.status,
		Priority:      c.priority,
		HeritageValue: float64(c.heritageValue) / 100,
	}
}

type CreateClientDto struct {
	Name       string  `json:"cliente_nome" validate:"required"`
	Email      string  `json:"cliente_email" validate:"required,email"`
	ClientType string  `json:"tipo_solicitacao" validate:"required"`
	Value      float64 `json:"valor_patrimonio" validate:"required,gt=0"`
}
