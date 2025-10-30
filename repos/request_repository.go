package repos

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
	"time"
)

type RequestRepository struct {
	db db.ISqlExecutor
}

func NewRequestRepository(db db.ISqlExecutor) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Add(ctx context.Context, request domain.Request) error {
	query := `INSERT INTO requests (id, type, payload, status, source_id, source_name) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, request.ID, request.Type, request.Payload, request.Status, request.SourceID, request.SourceName)
	return err
}

func (r *RequestRepository) Remove(ctx context.Context, rID uuid.UUID) error {
	query := `UPDATE requests SET is_deleted = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, rID)
	return err
}

func (r *RequestRepository) Accept(ctx context.Context, rID uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE requests SET 
                    done = TRUE, 
                    edit_user_id = $1,
                    status = $2,
                    solved_at = $3
                WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, userID, serializable.StatusAccepted, time.Now(), rID)
	return err
}
func (r *RequestRepository) Decline(ctx context.Context, rID uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE requests SET
			 		status = $1,
			 		edit_user_id = $2,
			 		solved_at = $3,
			 		done = FALSE
				WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, serializable.StatusDeclined, userID, time.Now(), rID)
	return err
}
func (r *RequestRepository) Get(ctx context.Context, rID uuid.UUID) (*domain.Request, error) {
	var result domain.Request
	query := `SELECT 
    			id, 
    			type, 
    			payload, 
    			status, 
    			source_id,
    			source_name,
    			solved_at, 
    			edit_user_id, 
    			solved_at, 
    			done, 
    			rollback_at 
			FROM requests 
			WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, rID).Scan(
		&result.ID,
		&result.Type,
		&result.Payload,
		&result.Status,
		&result.SourceID,
		&result.SourceName,
		&result.SolvedAt,
		&result.EditUserID,
		&result.SolvedAt,
		&result.Done,
		&result.RollbackAt,
	)
	return &result, err
}

func (r *RequestRepository) ExistsBySourceIDAndType(ctx context.Context, srcID uuid.UUID, rType serializable.RequestType) (bool, error) {
	query := `SELECT id FROM requests WHERE source_id = $1 AND type = $2`
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, srcID, rType).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RequestRepository) WithTx(tx *sql.Tx) repos2.IRequestRepository {
	return &RequestRepository{db: tx}
}
