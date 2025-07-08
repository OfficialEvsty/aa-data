package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
)

type ItemRepository struct {
	exec db.ISqlExecutor
}

func NewItemRepository(exec db.ISqlExecutor) *ItemRepository {
	return &ItemRepository{exec}
}

func (r *ItemRepository) WithTx(tx *sql.Tx) repos.IItemRepository {
	return &ItemRepository{tx}
}

func (r *ItemRepository) Add(
	ctx context.Context,
	temp domain.AAItemTemplate,
) (*domain.AAItemTemplate, error) {
	var result domain.AAItemTemplate
	query := `INSERT INTO aa_items (id, name, tier, img_grade_url, img_url)
  			  VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id 
			  RETURNING id, name, tier, img_grade_url, img_url;`
	row := r.exec.QueryRowContext(ctx, query, temp.ID, temp.Name, temp.Tier, temp.TierURL, temp.ImageURL)
	err := row.Scan(result.ID, result.Name, result.Tier, result.ImageURL, result.ImageURL)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ItemRepository) Remove(ctx context.Context, id int64) error {
	return nil
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*domain.AAItemTemplate, error) {
	var result domain.AAItemTemplate
	return &result, nil
}

func (r *ItemRepository) List(ctx context.Context) ([]*domain.AAItemTemplate, error) {
	var result []*domain.AAItemTemplate
	return result, nil
}
