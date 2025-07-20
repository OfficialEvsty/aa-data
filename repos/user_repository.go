package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
	"time"
)

type UserRepository struct {
	exec db.ISqlExecutor
}

func NewUserRepository(executor db.ISqlExecutor) *UserRepository {
	return &UserRepository{
		exec: executor,
	}
}

func (r *UserRepository) AddOrUpdate(ctx context.Context, user domain.User) (*domain.User, error) {
	updatedUserActivity := time.Now()
	query := `INSERT INTO users (id, username, email)
			  VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET username = $2, email = $3, last_seen = $4
			  RETURNING id, username, email, created_at, last_seen`
	row := r.exec.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, updatedUserActivity)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastSeen)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, created_at, last_seen FROM users WHERE id = $1`
	row := r.exec.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastSeen)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, created_at, last_seen FROM users WHERE email = $1`
	row := r.exec.QueryRowContext(ctx, query, email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.LastSeen)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//func (r *UserRepository) Update(ctx context.Context, user domain.User) (*domain.User, error) {
//	var updatedUser domain.User
//	query := `UPDATE users
//			  SET username = $2, email = $3, last_seen = $4
//			  WHERE id = $1 RETURNING id, username, email, created_at, last_seen`
//	row := r.exec.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, time.Now())
//	err := row.Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.LastSeen)
//	if err != nil {
//		return nil, err
//	}
//	return &updatedUser, nil
//}

func (r *UserRepository) Remove(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) WithTx(tx *sql.Tx) repos.IUserRepository {
	return &UserRepository{exec: tx}
}
