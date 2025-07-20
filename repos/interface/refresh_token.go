package repos

import (
	"context"
	"database/sql"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/google/uuid"
)

// IRefreshTokenRepository crud operations under refresh token data
type IRefreshTokenRepository interface {
	AddOrUpdate(context.Context, domain.RefreshToken) (*domain.RefreshToken, error)
	GetByUserID(context.Context, uuid.UUID) (*domain.RefreshToken, error)
	GetByToken(context.Context, string) (*domain.RefreshToken, error)
	Remove(context.Context, string) error
	WithTx(*sql.Tx) IRefreshTokenRepository
}
