package input

import (
	"github.com/efritz/lunar-fever/internal/common/datastructures"
	"github.com/go-gl/glfw/v3.2/glfw"
	"golang.org/x/exp/maps"
)

type KeyboardState struct {
	curr KeySet
	prev KeySet
	next KeySet
}

type KeySet datastructures.Set[glfw.Key]

func NewKeyboardState() *KeyboardState {
	return &KeyboardState{
		curr: KeySet{},
		prev: KeySet{},
		next: KeySet{},
	}
}

func (s *KeyboardState) OnKeyChange(key glfw.Key, action glfw.Action) {
	switch action {
	case glfw.Press:
		s.next[key] = struct{}{}

	case glfw.Release:
		delete(s.next, key)
	}
}

func (s *KeyboardState) Update() {
	s.curr, s.prev, s.next = s.next, s.curr, s.prev
	maps.Clear(s.next)
	maps.Copy(s.next, s.curr)
}

func (s *KeyboardState) IsKeyDown(key glfw.Key) bool {
	_, ok := s.curr[key]
	return ok
}

func (s *KeyboardState) WasKeyDown(key glfw.Key) bool {
	_, ok := s.prev[key]
	return ok
}

func (s *KeyboardState) IsKeyNewlyDown(key glfw.Key) bool {
	return s.IsKeyDown(key) && !s.WasKeyDown(key)
}
