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

type ChainRepository struct {
	exec db.ISqlExecutor
}

func NewChainRepository(exec db.ISqlExecutor) *ChainRepository {
	return &ChainRepository{exec}
}

func (r *ChainRepository) Add(ctx context.Context, chain domain.NicknameChain) error {
	query := `INSERT INTO chains (chain_id, parent_chain_id, nickname_id)
			  VALUES ($1, $2, $3) ON CONFLICT (chain_id) DO NOTHING`
	_, err := r.exec.ExecContext(ctx, query, chain.ChainID, chain.ParentChainID, chain.NicknameID)
	return err
}

// GetChain receive all chain ids []*domain.NicknameChain by provided nicknameID
func (r *ChainRepository) GetChain(ctx context.Context, nicknameID uuid.UUID) ([]*domain.NicknameChain, error) {
	var result []*domain.NicknameChain
	query := `SELECT chain_id, parent_chain_id, nickname_id, created_at, active, chained_at
			  FROM chains 
			  WHERE nickname_id = $1 
			  ORDER BY chained_at DESC`
	rows, err := r.exec.QueryContext(ctx, query, nicknameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var chain domain.NicknameChain
		err = rows.Scan(
			&chain.ChainID,
			&chain.ParentChainID,
			&chain.NicknameID,
			&chain.CreatedAt,
			&chain.Active,
			&chain.ChainedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &chain)
	}
	return result, err
}

func (r *ChainRepository) GetRootChain(ctx context.Context, chainID uuid.UUID) (*domain.NicknameChain, error) {
	var result domain.NicknameChain
	query := `WITH RECURSIVE root_chain AS (
				SELECT chain_id, nickname_id, parent_chain_id, created_at, active, chained_at
				FROM chains
				WHERE chain_id = $1  -- стартуем с любого узла
			
				UNION ALL
			
				SELECT c.chain_id, c.nickname_id, c.parent_chain_id, c.created_at, c.active, c.chained_at
				FROM chains c
				INNER JOIN root_chain rc ON rc.parent_chain_id = c.chain_id
			)
			SELECT chain_id, nickname_id, parent_chain_id, active, created_at, chained_at
			FROM root_chain
			WHERE parent_chain_id IS NULL; -- находим корень`
	err := r.exec.QueryRowContext(ctx, query, chainID).Scan(&result.ChainID, &result.NicknameID, &result.ParentChainID, &result.Active, &result.CreatedAt, &result.ChainedAt)
	return &result, err
}

func (r *ChainRepository) GetActiveChainID(ctx context.Context, chainID uuid.UUID) (*domain.NicknameChain, error) {
	var result domain.NicknameChain
	query := `WITH RECURSIVE chain_tree AS (
			SELECT chain_id, nickname_id, parent_chain_id, created_at, active, chained_at
			FROM chains
			WHERE chain_id = $1

			UNION ALL
			
			SELECT c.chain_id, c.nickname_id, c.parent_chain_id, c.created_at, c.active, c.chained_at
			FROM chains c
			INNER JOIN chain_tree ct ON c.parent_chain_id = ct.chain_id
		)
		SELECT DISTINCT chain_id, nickname_id, parent_chain_id, created_at, active, chained_at
		FROM chain_tree
		WHERE active = TRUE
		LIMIT 1;`
	err := r.exec.QueryRowContext(ctx, query, chainID).Scan(&result.ChainID, &result.NicknameID, &result.ParentChainID, &result.CreatedAt, &result.Active, &result.ChainedAt)
	return &result, err
}

func (r *ChainRepository) GetNicknameChains(ctx context.Context, chainID uuid.UUID) ([]*domain.NicknameChain, error) {
	var result []*domain.NicknameChain
	query := `WITH RECURSIVE chain_tree AS (
            SELECT chain_id, nickname_id, parent_chain_id, created_at, active, chained_at
            FROM chains
            WHERE chain_id = $1

            UNION ALL

            SELECT c.chain_id, c.nickname_id, c.parent_chain_id, c.created_at, c.active, c.chained_at
            FROM chains c
            INNER JOIN chain_tree ct ON c.parent_chain_id = ct.chain_id
        )
        SELECT DISTINCT chain_id, nickname_id, parent_chain_id, chained_at, active FROM chain_tree;`
	rows, err := r.exec.QueryContext(ctx, query, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var chain domain.NicknameChain
		err = rows.Scan(&chain.ChainID, &chain.NicknameID, &chain.ParentChainID, &chain.ChainedAt, &chain.Active)
		if err != nil {
			return nil, err
		}
		result = append(result, &chain)
	}
	return result, nil
}
func (r *ChainRepository) AttachChain(ctx context.Context, parent uuid.UUID, child uuid.UUID) error {
	query := `UPDATE chains
   			  SET parent_chain_id = $1
    		  WHERE chain_id = $2;`
	_, err := r.exec.ExecContext(ctx, query, parent, child, time.Now())
	return err
}

// DetachChain detaches child from specified parent
func (r *ChainRepository) DetachChain(ctx context.Context, parent uuid.UUID) error {
	query := `UPDATE chains SET parent_chain_id = NULL WHERE parent_chain_id = $1;`
	_, err := r.exec.ExecContext(ctx, query, parent)
	return err
}

func (r *ChainRepository) UpdateStatus(ctx context.Context, chainID uuid.UUID, active bool) error {
	query := `UPDATE chains
    		  SET active = $2
    		  WHERE chain_id = $1;`
	_, err := r.exec.ExecContext(ctx, query, chainID, active)
	return err
}

func (r *ChainRepository) UpdateChainedAt(ctx context.Context, parent uuid.UUID, chainedAt *time.Time) error {
	query := `UPDATE chains
			  SET chained_at = $1
			  WHERE chain_id = $2;`
	_, err := r.exec.ExecContext(ctx, query, chainedAt, parent)
	return err
}

func (r *ChainRepository) WithTx(tx *sql.Tx) repos.IChainRepository {
	return &ChainRepository{tx}
}
