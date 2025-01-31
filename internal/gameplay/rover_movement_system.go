package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type roverMovementSystem struct {
	*GameContext
}

func NewRoverMovementSystem(ctx *GameContext) system.System {
	return &roverMovementSystem{GameContext: ctx}
}

func (s *roverMovementSystem) Init() {}
func (s *roverMovementSystem) Exit() {}

// TODO - deglobalize
var controllingRover = false

func (g *roverMovementSystem) Process(elapsedMs int64) {
	if g.Keyboard.IsKeyNewlyDown(glfw.KeyR) {
		controllingRover = !controllingRover
	}
	if !controllingRover {
		return
	}

	roverXDir := int64(0)
	roverYDir := int64(0)
	if g.Keyboard.IsKeyDown(glfw.KeyS) {
		roverYDir++
	}
	if g.Keyboard.IsKeyDown(glfw.KeyW) {
		roverYDir--
	}
	if g.Keyboard.IsKeyDown(glfw.KeyD) {
		roverXDir++
	}
	if g.Keyboard.IsKeyDown(glfw.KeyA) {
		roverXDir--
	}

	for _, entity := range g.RoverCollection.Entities() {
		component, ok := g.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		// angle := math.Atan232(my-component.Body.Position.Y, mx-component.Body.Position.X)
		// if angle < 0 {
		// 	angle = (2 * stdmath.Pi) - (-angle)
		// }
		//
		// 	component.Body.SetOrient(angle)

		mod := float32(1000)
		if roverXDir != 0 {
			dx := stdmath.Pi * float32(roverXDir) / (128 + 64)
			tireRotation, _ = math.Clamp(tireRotation+dx*2, -stdmath.Pi/6, stdmath.Pi/6)

			if roverYDir != 0 {
				component.Body.SetOrient(component.Body.Orient + dx*-float32(roverYDir))
			}
		} else {
			if tireRotation > 0 {
				tireRotation -= stdmath.Pi / 128
			} else {
				tireRotation += stdmath.Pi / 128
			}
		}

		if roverYDir != 0 {
			speed := float32(.35)
			transitionSpeed := float32(4)

			component.Body.LinearVelocity = component.Body.LinearVelocity.
				Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).
				Add(component.Body.Rotation.Mul(math.Vector{0, float32(roverYDir)}).
					Muls(speed * float32(elapsedMs) / mod * transitionSpeed))
		} else {
			transitionSpeed := float32(8)

			component.Body.LinearVelocity = component.Body.LinearVelocity.Muls(1 - (float32(elapsedMs) / mod * transitionSpeed))
			component.Body.AngularVelocity = component.Body.AngularVelocity * (1 - (float32(elapsedMs) / mod * transitionSpeed))
		}
	}
}
