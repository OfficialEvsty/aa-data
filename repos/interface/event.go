package repos

import (
	"context"
	"github.com/OfficialEvsty/aa-data/domain"
)

type IEventRepository interface {
	Add(context.Context, domain.Event) (*domain.Event, error)
	GetByID(context.Context, uint64) (*domain.Event, error)
	Remove(context.Context, uint64) error
}
