package repos

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	repos "github.com/OfficialEvsty/aa-data/repos/interface"
	"github.com/lib/pq"
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
	err := row.Scan(&result.ID, &result.Name, &result.Tier, &result.ImageURL, &result.ImageURL)
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

func (r *ItemRepository) GetByIDs(ctx context.Context, ids []int64) ([]*domain.AAItemTemplate, error) {
	query := `SELECT id, name, tier, img_grade_url, img_url
              FROM aa_items
              WHERE id = ANY($1)`
	rows, err := r.exec.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*domain.AAItemTemplate, 0)
	for rows.Next() {
		var item domain.AAItemTemplate
		err = rows.Scan(&item.ID, &item.Name, &item.Tier, &item.ImageURL, &item.ImageURL)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *ItemRepository) List(ctx context.Context) ([]*domain.AAItemTemplate, error) {
	var result []*domain.AAItemTemplate
	return result, nil
}
