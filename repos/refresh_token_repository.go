package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
)

type RefreshTokenRepository struct {
	exec db.ISqlExecutor
}

func NewRefreshTokenRepository(exec db.ISqlExecutor) *RefreshTokenRepository {
	return &RefreshTokenRepository{exec}
}

// AddOrUpdate adds or updates refresh token
func (r *RefreshTokenRepository) AddOrUpdate(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error) {
	query := `INSERT INTO tokens (token, user_id, expires_at)
		 	  VALUES ($1, $2, $3) ON CONFLICT (user_id)
			  DO UPDATE SET token = $1, expires_at = $3
			  RETURNING token, user_id, expires_at`
	row := r.exec.QueryRowContext(ctx, query, token.Token, token.UserID, token.ExpiresAt)
	err := row.Scan(&token.Token, &token.UserID, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var res domain.RefreshToken
	query := `SELECT * FROM tokens WHERE token = $1`
	row := r.exec.QueryRowContext(ctx, query, token)
	err := row.Scan(&res.Token, &res.UserID, &res.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.RefreshToken, error) {
	var res domain.RefreshToken
	query := `SELECT * FROM tokens WHERE user_id = $1`
	row := r.exec.QueryRowContext(ctx, query, userID)
	err := row.Scan(&res.Token, &res.UserID, &res.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *RefreshTokenRepository) Remove(ctx context.Context, token string) error {
	query := `DELETE FROM tokens WHERE token = $1`
	_, err := r.exec.ExecContext(ctx, query, token)
	return err
}

func (r *RefreshTokenRepository) WithTx(tx *sql.Tx) repos.IRefreshTokenRepository {
	return &RefreshTokenRepository{exec: tx}
}
