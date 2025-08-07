package commands

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/google/uuid"
)

type DeleteIncompleteRaidByUserIdCommand struct {
	UserID uuid.UUID
	RaidID uuid.UUID
}

type IncompleteRaidCleaner struct {
	exec db.ISqlExecutor
}

func NewIncompleteRaidCleaner(db db.ISqlExecutor) *IncompleteRaidCleaner {
	return &IncompleteRaidCleaner{
		exec: db,
	}
}

func (cleaner *IncompleteRaidCleaner) Handle(ctx context.Context, cmd *DeleteIncompleteRaidByUserIdCommand) error {
	query := `UPDATE raids r SET is_deleted = TRUE 
              FROM tenant_publishes tp 
              WHERE r.publish_id = tp.publish_id AND r.is_deleted = FALSE AND tp.user_id = $1 AND r.id = $2`
	_, err := cleaner.exec.ExecContext(ctx, query, cmd.UserID, cmd.RaidID)
	if err != nil {
		return err
	}
	return nil
}
