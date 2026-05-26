package clients

type ClientRepository struct{}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (r *ClientRepository) CreateClient(client *Client) (*Client, error) {
	return client, nil
}
