package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/google/uuid"
	"time"
)

type RaidRepository struct {
	exec db.ISqlExecutor
}

func NewRaidRepository(exec db.ISqlExecutor) *RaidRepository {
	return &RaidRepository{exec}
}

func (r *RaidRepository) Add(ctx context.Context, raid domain.Raid) error {
	query := `INSERT INTO raids (id, publish_id, status) VALUES ($1, $2, $3)`
	_, err := r.exec.ExecContext(ctx, query, raid.ID, raid.PublishID, raid.Status)
	if err != nil {
		return err
	}
	return nil
}
func (r *RaidRepository) Update(ctx context.Context, raid domain.Raid) error {
	query := `UPDATE raids SET publish_id=$2, raid_at=$3, attendance=$4, status=$5 WHERE id=$1`
	_, err := r.exec.ExecContext(ctx, query, raid.ID, raid.PublishID, raid.RaidAt, raid.Attendance, raid.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *RaidRepository) UpdateTiming(ctx context.Context, id uuid.UUID, raidAt time.Time) error {
	query := `UPDATE raids SET raid_at = $2 WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, id, raidAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *RaidRepository) UpdateAttendance(ctx context.Context, id uuid.UUID, attendance int) error {
	query := `UPDATE raids SET attendance = $2 WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, id, attendance)
	if err != nil {
		return err
	}
	return nil
}
func (r *RaidRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status serializable.Status) error {
	query := `UPDATE raids SET status = $2 WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, id, status)
	if err != nil {
		return err
	}
	return nil
}

func (r *RaidRepository) UpdateEndDateAndStatus(
	ctx context.Context,
	id uuid.UUID,
	raidAt time.Time,
	status serializable.Status,
) error {
	query := `UPDATE raids SET status = $2, raid_at = $3 WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, id, status, raidAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *RaidRepository) Remove(ctx context.Context, raidID uuid.UUID) error {
	query := `UPDATE raids SET is_deleted = TRUE WHERE id = $1`
	_, err := r.exec.ExecContext(ctx, query, raidID)
	if err != nil {
		return err
	}
	return nil
}
func (r *RaidRepository) GetById(ctx context.Context, raidID uuid.UUID) (*domain.Raid, error) {
	var raid domain.Raid
	query := `SELECT id, publish_id, raid_at, attendance, status, created_at FROM raids WHERE id = $1 AND is_deleted = FALSE`
	row := r.exec.QueryRowContext(ctx, query, raidID)
	err := row.Scan(
		&raid.ID,
		&raid.PublishID,
		&raid.RaidAt,
		&raid.Attendance,
		&raid.Status,
		&raid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &raid, nil
}
func (r *RaidRepository) WithTx(tx *sql.Tx) repos.IRaidRepository {
	return &RaidRepository{tx}
}
