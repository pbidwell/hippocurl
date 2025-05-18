/*
Copyright © 2025 Pablo Bidwell <bidwell.pablo@gmail.com>
*/
package modules

import (
	"github.com/pbidwell/hippocurl/internal/config"
)

type HippoModule interface {
	Name() string
	Description() string
	Logo() string
	Use() string
	Execute(config *config.App, args []string)
}
