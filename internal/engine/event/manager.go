package event

type Manager struct {
	listeners map[EventType][]any
}

func NewManager() *Manager {
	return &Manager{
		listeners: map[EventType][]any{},
	}
}
