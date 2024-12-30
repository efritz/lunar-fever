package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type playerMovementSystem struct {
	*GameContext
}

func NewPlayerMovementSystem(ctx *GameContext) system.System {
	return &playerMovementSystem{GameContext: ctx}
}

func (s *playerMovementSystem) Init() {
	NewEntityDeathEventManager(s.GameContext.EventManager).AddListener(s)
}

func (s *playerMovementSystem) Exit() {}

func (g *playerMovementSystem) Process(elapsedMs int64) {
	if controllingRover {
		return
	}

	playerXDir := float32(0)
	playerYDir := float32(0)
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
		playerXDir *= 1.5
		playerYDir *= 1.5
	}

	mx := g.Camera.Unprojectx(float32(g.Mouse.X()))
	my := g.Camera.UnprojectY(float32(g.Mouse.Y()))

	for _, entity := range g.PlayerCollection.Entities() {
		physicsComponent, ok := g.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		if healthComponent, ok := g.HealthComponentManager.GetComponent(entity); !ok || healthComponent.Health <= 0 {
			physicsComponent.Body.LinearVelocity = math.Vector{0, 0}
			physicsComponent.Body.AngularVelocity = 0
			continue
		}

		angle := math.Atan232(my-physicsComponent.Body.Position.Y, mx-physicsComponent.Body.Position.X)
		if angle < 0 {
			angle = (2 * stdmath.Pi) - (-angle)
		}
		angle -= float32(stdmath.Pi / 2)

		if physicsComponent.Body.Orient != angle {
			physicsComponent.Body.SetOrient(angle)
		}

		mod := float32(1000)
		if playerXDir != 0 || playerYDir != 0 {
			speed := float32(.35)
			transitionSpeed := float32(4)

			physicsComponent.Body.LinearVelocity = physicsComponent.Body.LinearVelocity.
				Muls(1 - (float32(elapsedMs) / mod * transitionSpeed)).
				Add(math.Vector{playerXDir, playerYDir}.Muls(speed * float32(elapsedMs) / mod * transitionSpeed))
		} else {
			transitionSpeed := float32(8)

			physicsComponent.Body.LinearVelocity = physicsComponent.Body.LinearVelocity.Muls(1 - (float32(elapsedMs) / mod * transitionSpeed))
			physicsComponent.Body.AngularVelocity = physicsComponent.Body.AngularVelocity * (1 - (float32(elapsedMs) / mod * transitionSpeed))
		}
	}
}

func (s *playerMovementSystem) OnEntityDeath(e EntityDeathEvent) {
	// s.physicsComponentManager.RemoveComponent(e.entity)
}
