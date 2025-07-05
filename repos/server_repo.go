package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type ServerRepository struct {
	exec db.ISqlExecutor
}

func NewServerRepository(exec db.ISqlExecutor) *ServerRepository {
	return &ServerRepository{exec}
}

// Add domain.AAServer to database
func (r *ServerRepository) Add(ctx context.Context, server domain.AAServer) (*domain.AAServer, error) {
	var result domain.AAServer
	query := `INSERT INTO aa_servers (id, name, external_id)
 			  VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id, name, external_id;`
	res := r.exec.QueryRowContext(ctx, query, server.ID, server.Name, server.ExternalID)
	err := res.Scan(&result.ID, &result.Name, &result.ExternalID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByExternalID for working with official RU Archeage server's ids
func (r *ServerRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.AAServer, error) {
	var result domain.AAServer
	query := `SELECT id, name, external_id FROM aa_servers WHERE external_id = $1`
	res := r.exec.QueryRowContext(ctx, query, externalID)
	err := res.Scan(&result.ID, &result.Name, &result.ExternalID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (r *ServerRepository) List(ctx context.Context) ([]*domain.AAServer, error) {
	var result []*domain.AAServer
	query := `SELECT id, name, external_id FROM aa_servers`
	rows, err := r.exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var server domain.AAServer
		if err = rows.Scan(&server.ID, &server.Name, &server.ExternalID); err != nil {
			return nil, err
		}
		result = append(result, &server)
	}
	return result, nil
}
func (r *ServerRepository) Remove(context.Context, uuid.UUID) error {
	return nil
}

func (r *ServerRepository) WithTx(tx *sql.Tx) repos.IServerRepository {
	if r == nil {
		panic("ServerRepository.WithTx called on nil receiver")
	}
	if tx == nil {
		panic("ServerRepository.WithTx received nil tx")
	}
	return &ServerRepository{tx}
}
