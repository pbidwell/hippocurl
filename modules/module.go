package modules

type HippoModule interface {
	Name() string
	Description() string
	Logo() string
	Execute(args []string)
}
