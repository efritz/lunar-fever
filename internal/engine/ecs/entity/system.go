package entity

import "github.com/efritz/lunar-fever/internal/engine/ecs/system"

type System struct {
	delegate   SystemDelegate
	collection *Collection
}

type SystemDelegate interface {
	Init()
	Exit()
	Process(entity Entity, elapsedMs int64)
}

func NewSystem(delegate SystemDelegate, collection *Collection) system.System {
	return &System{
		delegate:   delegate,
		collection: collection,
	}
}

func (s *System) Init() {
	s.delegate.Init()
}

func (s *System) Exit() {
	s.delegate.Exit()
}

func (s *System) Process(elapsedMs int64) {
	for _, entity := range s.collection.Entities() {
		s.delegate.Process(entity, elapsedMs)
	}
}
