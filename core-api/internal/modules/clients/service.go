package clients

type ClientService struct {
	repository *ClientRepository
}

func NewClientService() *ClientService {
	repository := NewClientRepository()
	return &ClientService{
		repository: repository,
	}
}

func (s *ClientService) CreateClient(client *CreateClientDto) (*Client, error) {
	return s.repository.CreateClient(&Client{
		name:          client.Name,
		email:         client.Email,
		clientType:    client.ClientType,
		heritageValue: client.Value,
	})
}
