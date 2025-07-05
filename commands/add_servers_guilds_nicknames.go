package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	repos2 "github.com/OfficialEvsty/aa-data/repos/interface"
)

// AddServersGuildsNicknamesCommand add all related entities together
type AddServersGuildsNicknamesCommand struct {
	Servers []*usecase.ServerData
}

type ServerImporter struct {
	tx         db.ITxExecutor
	serverRepo repos2.IServerRepository
	guildRepo  repos2.IGuildRepository
	nickRepo   repos2.INicknameRepository
	linkRepo   repos2.IGuildNicknameRepository
}

func NewServerImporter(
	tx db.ITxExecutor,
	serverRepo repos2.IServerRepository,
	guildRepo repos2.IGuildRepository,
	nickRepo repos2.INicknameRepository,
	linkRepo repos2.IGuildNicknameRepository,
) *ServerImporter {
	return &ServerImporter{
		tx:         tx,
		serverRepo: serverRepo,
		guildRepo:  guildRepo,
		nickRepo:   nickRepo,
		linkRepo:   linkRepo,
	}
}

// Handle executes command
func (si *ServerImporter) Handle(ctx context.Context, cmd AddServersGuildsNicknamesCommand) error {
	err := si.tx.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		for _, s := range cmd.Servers {
			var server *domain.AAServer
			server, err := si.serverRepo.WithTx(tx).GetByExternalID(ctx, s.Server.ExternalID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					server, err = si.serverRepo.WithTx(tx).Add(ctx, *s.Server)
					if err != nil {
						return fmt.Errorf("error adding server: %w", err)
					}
				} else {
					return fmt.Errorf("error getting server by external ID: %w", err)
				}
			}
			for _, g := range s.Guilds {
				g.Guild.ServerID = server.ID
				guild, err := si.guildRepo.WithTx(tx).Add(ctx, *g.Guild)
				if err != nil {
					return fmt.Errorf("error adding guild: %v", err)
				}
				for _, n := range g.Nicknames {
					n.Nickname.ServerID = server.ID
					nickname, err := si.nickRepo.WithTx(tx).Create(ctx, n.Nickname)
					if err != nil {
						return fmt.Errorf("error creating nickname: %v", err)
					}
					err = si.linkRepo.WithTx(tx).Add(ctx, guild.ID, nickname.ID)
					if err != nil {
						return fmt.Errorf("error adding nickname: %v", err)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
