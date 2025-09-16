package queries

import (
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/google/uuid"
)

type ChainedNickname struct {
	ChainID    uuid.UUID `json:"chain_id"`
	NicknameID uuid.UUID `json:"nickname_id"`
	GuildID    uuid.UUID `json:"guild_id"`
	Name       string    `json:"name"`
	GuildName  string    `json:"guild_name"`
}

type GetChainedParticipantsByTenantIDQuery struct {
	exec db.ISqlExecutor
}

func NewGetChainedParticipantsByTenantIDQuery(exec db.ISqlExecutor) *GetChainedParticipantsByTenantIDQuery {
	return &GetChainedParticipantsByTenantIDQuery{exec: exec}
}

//func (q *GetChainedParticipantsByTenantIDQuery) Handle(ctx context.Context, chainIDs []uuid.UUID) ([]*ChainedNickname, error) {
//	query := `
//				SELECT
//			 `
//}
