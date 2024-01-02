package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type MouseState struct {
	curr FrozenMouseState
	prev FrozenMouseState
	next FrozenMouseState
}

type FrozenMouseState struct {
	X           float64
	Y           float64
	ScrollWheel float64
	LeftButton  bool
	RightButton bool
}

func NewMouseState() *MouseState {
	return &MouseState{}
}

func (s *MouseState) OnPositionChange(xpos, ypos float64) {
	s.next.X = xpos
	s.next.Y = ypos
}

func (s *MouseState) OnScrollChange(yoff float64) {
	s.next.ScrollWheel += yoff
}

func (s *MouseState) OnButtonChange(button glfw.MouseButton, action glfw.Action) {
	switch button {
	case glfw.MouseButtonLeft:
		s.next.LeftButton = action == glfw.Press

	case glfw.MouseButtonRight:
		s.next.RightButton = action == glfw.Press
	}
}

func (s *MouseState) Update() {
	s.curr, s.prev = s.next, s.curr
}

func (s *MouseState) X() float64                 { return s.curr.X }
func (s *MouseState) Y() float64                 { return s.curr.Y }
func (s *MouseState) LeftButton() bool           { return s.curr.LeftButton }
func (s *MouseState) LeftButtonNewlyDown() bool  { return s.curr.LeftButton && !s.prev.LeftButton }
func (s *MouseState) RightButton() bool          { return s.curr.RightButton }
func (s *MouseState) RightButtonNewlyDown() bool { return s.curr.RightButton && !s.prev.RightButton }
func (s *MouseState) ScrollDelta() float64       { return s.curr.ScrollWheel - s.prev.ScrollWheel }
