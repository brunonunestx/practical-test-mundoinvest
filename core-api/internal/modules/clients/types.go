package clients

type Client struct {
	id            int
	name          string
	email         string
	clientType    string
	status        string
	heritageValue float64
}

type CreateClientDto struct {
	Name       string  `json:"cliente_nome"`
	Email      string  `json:"cliente_email"`
	ClientType string  `json:"tipo_cliente"`
	Value      float64 `json:"valor_patrimonio"`
}
