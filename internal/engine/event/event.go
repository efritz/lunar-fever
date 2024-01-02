package event

type EventType interface{}

type Event[T EventType, L any] interface {
	EventType() T
	Notify(L)
}
