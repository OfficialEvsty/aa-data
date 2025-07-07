package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
)

// BossesRepository provide crud operations under aa-bosses table
type BossesRepository struct {
	exec db.ISqlExecutor
}

func NewBossesRepository(exec db.ISqlExecutor) *BossesRepository {
	return &BossesRepository{exec}
}

// WithTx switch modes between tx and sql.db
func (r *BossesRepository) WithTx(tx *sql.Tx) repos2.IBossesRepository {
	return &BossesRepository{tx}
}

// Add implementation of adding boss to table
func (r *BossesRepository) Add(ctx context.Context, boss domain.AABoss) (*domain.AABoss, error) {
	var result domain.AABoss
	query := `INSERT INTO aa_bosses (id, name, drop, level, img_url)
			  VALUES ($1, $2, $3, $4, $4) ON CONFLICT DO UPDATE SET drop = $3
			  RETURNING id, name, drop, level, img_url`
	err := r.exec.QueryRowContext(ctx, query, boss.ID, boss.Name, boss.Loot, boss.Level, boss.ImageURL).Scan(
		&result.ID,
		&result.Name,
		&result.Loot,
		&result.Level,
		&result.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Remove implementation of removing boss record from table aa-bosses
func (r *BossesRepository) Remove(ctx context.Context, id int64) error {
	return nil
}

func (r *BossesRepository) GetByID(ctx context.Context, id int64) (*domain.AABoss, error) {
	var boss domain.AABoss
	query := `SELECT id, name, drop, level, img_url  FROM aa_bosses WHERE id = $1`
	err := r.exec.QueryRowContext(ctx, query, id).Scan(
		&boss.ID,
		&boss.Name,
		&boss.Loot,
		&boss.Level,
		&boss.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	return &boss, nil
}

func (r *BossesRepository) List(ctx context.Context) ([]*domain.AABoss, error) {
	var bosses []*domain.AABoss
	return bosses, nil
}
