package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type cameraMovementSystem struct {
	*engine.Context
}

func (s *cameraMovementSystem) Init() {}
func (s *cameraMovementSystem) Exit() {}

func (s *cameraMovementSystem) Process(elapsedMs int64) {
	cameraXDir := int64(0)
	cameraYDir := int64(0)
	if s.Keyboard.IsKeyDown(glfw.KeyUp) {
		cameraYDir++
	}
	if s.Keyboard.IsKeyDown(glfw.KeyDown) {
		cameraYDir--
	}
	if s.Keyboard.IsKeyDown(glfw.KeyLeft) {
		cameraXDir++
	}
	if s.Keyboard.IsKeyDown(glfw.KeyRight) {
		cameraXDir--
	}

	mod := float32(500) // TODO - why so slow?
	s.Camera.Translate(float32(cameraXDir*elapsedMs), float32(cameraYDir*elapsedMs))
	s.Camera.Zoom(float32(s.Mouse.ScrollDelta()) / mod)

}
