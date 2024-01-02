package component

import (
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/event"
)

type TypedManager[C Component[T], T ComponentType] struct {
	manager *Manager
}

func NewTypedManager[C Component[T], T ComponentType](manager *Manager, eventManager *event.Manager) *TypedManager[C, T] {
	return &TypedManager[C, T]{
		manager: manager,
	}
}

func (m *TypedManager[C, T]) HasComponent(e entity.Entity) bool {
	var componentType T // Infer component type value from type param
	return m.manager.HasComponent(e, componentType)
}

func (m *TypedManager[C, T]) GetComponent(e entity.Entity) (component C, _ bool) {
	var componentType T // Infer component type value from type param
	rawComponent, ok := m.manager.GetComponent(e, componentType)
	if !ok {
		return
	}

	component, ok = rawComponent.(C)
	if !ok {
		panic("malformed component manager")
	}

	return component, true
}

func (m *TypedManager[C, T]) AddComponent(e entity.Entity, component Component[T]) {
	var componentType T // Infer component type value from type param
	m.manager.AddComponent(e, component.(Component[any]), componentType)
}

func (m *TypedManager[C, T]) RemoveComponent(e entity.Entity) {
	var componentType T // Infer component type value from type param
	m.manager.RemoveComponent(e, componentType)
}
