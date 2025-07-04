package usecase

import "github.com/OfficialEvsty/aa-data/domain"

// ServerData provides data about server's guilds
type ServerData struct {
	Server *domain.AAServer
	Guilds []*GuildData
}

// GuildData provides data about guild's nicknames
type GuildData struct {
	Guild     *domain.AAGuild
	Nicknames []*domain.GuildNickname
}
