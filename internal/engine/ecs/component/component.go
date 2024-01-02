package component

type ComponentType interface{}

type Component[T ComponentType] interface {
	ComponentType() T
}
