package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type playerMovementSystem struct {
	*engine.Context
	playerCollection        *entity.Collection
	physicsComponentManager *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
}

func (s *playerMovementSystem) Init() {}
func (s *playerMovementSystem) Exit() {}

func (g *playerMovementSystem) Process(elapsedMs int64) {
	if controllingRover {
		return
	}

	playerXDir := int64(0)
	playerYDir := int64(0)
	if g.Keyboard.IsKeyDown(glfw.KeyS) {
		playerYDir++
	}
	if g.Keyboard.IsKeyDown(glfw.KeyW) {
		playerYDir--
	}
	if g.Keyboard.IsKeyDown(glfw.KeyD) {
		playerXDir++
	}
	if g.Keyboard.IsKeyDown(glfw.KeyA) {
		playerXDir--
	}

	if g.Keyboard.IsKeyDown(glfw.KeyLeftShift) {
		playerXDir *= 3
		playerYDir *= 3
	}

	mx := g.Camera.Unprojectx(float32(g.Mouse.X()))
	my := g.Camera.UnprojectY(float32(g.Mouse.Y()))

	for _, entity := range g.playerCollection.Entities() {
		component, ok := g.physicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		angle := math.Atan232(my-component.Body.Position.Y, mx-component.Body.Position.X)
		if angle < 0 {
			angle = (2 * stdmath.Pi) - (-angle)
		}

		component.Body.SetOrient(angle)

		mod := float32(100) // TODO - why so slow?
		if playerXDir != 0 || playerYDir != 0 {
			speed := float32(.35)
			transitionSpeed := float32(4)

			component.Body.LinearVelocity = component.Body.LinearVelocity.
				Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).
				Add(math.Vector{float32(playerXDir), float32(playerYDir)}).
				Muls(speed * float32(elapsedMs) / mod * transitionSpeed)
		} else {
			transitionSpeed := float32(8)

			component.Body.LinearVelocity = component.Body.LinearVelocity.Muls(1 - (float32(elapsedMs) / mod * transitionSpeed))
			component.Body.AngularVelocity = component.Body.AngularVelocity * (1 - (float32(elapsedMs) / mod * transitionSpeed))
		}
	}
}
