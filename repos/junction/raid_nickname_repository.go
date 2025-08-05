package junction_repos

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	junction_repos "github.com/OfficialEvsty/aa-data/repos/interface/junction"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
)

type RaidNicknameRepository struct {
	exec db.ISqlExecutor
}

func NewRaidNicknameRepository(exec db.ISqlExecutor) *RaidNicknameRepository {
	return &RaidNicknameRepository{exec}
}

func (r *RaidNicknameRepository) AddNicknames(ctx context.Context, raidID uuid.UUID, nicknameIDs []uuid.UUID) error {
	valueStrings := make([]string, 0, len(nicknameIDs))
	valueArgs := make([]interface{}, 0, len(nicknameIDs)*2)

	for i, nicknameID := range nicknameIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, raidID, nicknameID)
	}

	stmt := fmt.Sprintf("INSERT INTO attendance (raid_id, nickname_id) VALUES %s ON CONFLICT (raid_id, nickname_id) DO NOTHING", strings.Join(valueStrings, ","))
	_, err := r.exec.ExecContext(ctx, stmt, valueArgs...)
	return err
}
func (r *RaidNicknameRepository) RemoveNicknames(ctx context.Context, raidID uuid.UUID, nicknameIDs []uuid.UUID) error {
	query := `DELETE FROM attendance WHERE raid_id = $1 AND nickname_id = ANY($2)`
	_, err := r.exec.ExecContext(ctx, query, raidID, pq.Array(nicknameIDs))
	return err
}

func (r *RaidNicknameRepository) WithTx(tx *sql.Tx) junction_repos.IRaidNicknameRepository {
	return &RaidNicknameRepository{tx}
}
