package fabapi

type ActionType string

const (
	ActionAdd    ActionType = "ADD"
	ActionRemove ActionType = "REMOVE"
	ActionUpdate ActionType = "UPDATE"
)
