package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

// NicknameRepo nickname's repository implementation
type NicknameRepo struct {
	exec db.ISqlExecutor
}

// NewNicknameRepo creates instance of NicknameRepo
func NewNicknameRepo(executor db.ISqlExecutor) *NicknameRepo {
	return &NicknameRepo{
		exec: executor,
	}
}

// Create saves nickname in table aa_nicknames
func (r *NicknameRepo) Create(ctx context.Context, nickname domain.AANickname) (*domain.AANickname, error) {
	var result domain.AANickname

	query := `INSERT INTO aa_nicknames (id, server_id, name) 
			  VALUES ($1, $2, $3) 
			  ON CONFLICT (server_id, name) DO NOTHING 
			  RETURNING id, name, server_id, created_at`
	res := r.exec.QueryRowContext(ctx, query, uuid.New(), nickname.ServerID, nickname.Name)
	err := res.Scan(&result.ID, &result.Name, &result.ServerID, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByName returns domain.AANickname by his name
func (r *NicknameRepo) GetByName(ctx context.Context, serverID uuid.UUID, name string) (*domain.AANickname, error) {
	var result domain.AANickname
	query := `SELECT id, name, server_id, created_at
			  FROM aa_nicknames 
			  WHERE server_id = $1 AND name = $2`
	res, err := r.exec.QueryContext(ctx, query, serverID, name)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("deadline exceeded %v", err)
		}
		return nil, err
	}
	err = res.Scan(&result.ID, &result.Name, &result.ServerID, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *NicknameRepo) WithTx(tx *sql.Tx) repos.INicknameRepository {
	return &NicknameRepo{tx}
}
