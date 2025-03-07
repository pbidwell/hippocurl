package modules

import (
	"context"
)

type HippoModule interface {
	Name() string
	Description() string
	Logo() string
	Execute(ctx context.Context, args []string)
}
