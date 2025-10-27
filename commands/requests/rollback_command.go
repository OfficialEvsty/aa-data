package requests

import (
	"context"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
)

type RollbackCommand interface {
	Execute(context.Context) error
	Rollback(context.Context) error
	Type() serializable.RequestType
}
