package junction_repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
)

type LunarkRaidRepository struct {
	exec db.ISqlExecutor
}

func NewLunarkRaidRepository(exec db.ISqlExecutor) *LunarkRaidRepository {
	return &LunarkRaidRepository{exec}
}

func (r *LunarkRaidRepository) Add(ctx context.Context, lunark domain.LunarkRaid) error {
	query := `INSERT INTO lunark_raids (lunark_id, raid_id) VALUES ($1, $2)`
	_, err := r.exec.ExecContext(ctx, query, lunark.LunarkID, lunark.RaidID)
	if err != nil {
		return err
	}
	return nil
}
func (r *LunarkRaidRepository) Remove(ctx context.Context, raidID uuid.UUID) error {
	query := `DELETE FROM lunark_raids WHERE raid_id = $1`
	_, err := r.exec.ExecContext(ctx, query, raidID)
	if err != nil {
		return err
	}
	return nil
}
func (r *LunarkRaidRepository) All(ctx context.Context, lunarkID uuid.UUID) ([]*domain.LunarkRaid, error) {
	var result []*domain.LunarkRaid
	query := `SELECT lunark_id, raid_id FROM lunark_raids WHERE lunark_id = $1`
	rows, err := r.exec.QueryContext(ctx, query, lunarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var lr domain.LunarkRaid
		if err = rows.Scan(&lr.LunarkID, &lr.RaidID); err != nil {
			return nil, err
		}
		result = append(result, &lr)
	}
	return result, nil
}
func (r *LunarkRaidRepository) GetByID(ctx context.Context, raidID uuid.UUID) (*domain.LunarkRaid, error) {
	query := `SELECT lunark_id, raid_id FROM lunark_raids WHERE raid_id = $1`
	row := r.exec.QueryRowContext(ctx, query, raidID)
	var lr domain.LunarkRaid
	if err := row.Scan(&lr.LunarkID, &lr.RaidID); err != nil {
		return nil, err
	}
	return &lr, nil
}
func (r *LunarkRaidRepository) WithTx(tx *sql.Tx) junction_repos2.ILunarkRaidRepository {
	return &LunarkRaidRepository{tx}
}
