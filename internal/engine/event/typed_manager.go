package event

type TypedManager[E Event[T, L], T EventType, L any] struct {
	manager *Manager
}

func NewTypedManager[E Event[T, L], T EventType, L any](manager *Manager) *TypedManager[E, T, L] {
	return &TypedManager[E, T, L]{
		manager: manager,
	}
}

func (r *TypedManager[E, T, L]) AddListener(listener L) {
	var e E
	eventType := e.EventType() // Infer event type value from type param
	r.manager.listeners[eventType] = append(r.manager.listeners[eventType], listener)
}

func (r *TypedManager[E, T, L]) Dispatch(event E) {
	for _, listener := range r.manager.listeners[event.EventType()] {
		listenerL, ok := listener.(L)
		if !ok {
			panic("Listener type mismatch")
		}

		event.Notify(listenerL)
	}
}
