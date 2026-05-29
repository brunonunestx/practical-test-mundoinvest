package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	db "core-api/services/database/generated"
)

type ClientRepository struct {
	queries *db.Queries
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		queries: db.New(pool),
	}
}

func (r *ClientRepository) CreateClient(ctx context.Context, client *Client) (*Client, error) {
	row, err := r.queries.CreateClient(ctx, db.CreateClientParams{
		Name:        client.name,
		Email:       client.email,
		RequestType: client.clientType,
		Status:      db.RequestStatusEnumPENDINGANALYSIS,
		Priority:    db.PriorityEnum(client.priority),
		Amount:      int32(client.heritageValue),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		name:          row.Name,
		email:         row.Email,
		clientType:    row.RequestType,
		status:        string(row.Status),
		priority:      string(row.Priority),
		heritageValue: client.heritageValue,
	}, nil
}
