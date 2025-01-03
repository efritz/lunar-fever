package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type npcMovementRenderSystem struct {
	*GameContext
	emptyTexture rendering.Texture
}

func NewNpcMovementRenderSystem(ctx *GameContext) system.System {
	return &npcMovementRenderSystem{GameContext: ctx}
}

func (s *npcMovementRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (s *npcMovementRenderSystem) Exit() {}

var elapsed = int64(0)

func (s *npcMovementRenderSystem) Process(elapsedMs int64) {
	s.SpriteBatch.Begin()

	elapsed += elapsedMs
	debug := elapsed < 500
	elapsed = elapsed % 1000

	for _, entity := range s.NpcCollection.Entities() {
		// physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		// if !ok {
		// 	continue
		// }

		pathfindingComponent, ok := s.PathfindingComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		prev := math.Vector{}

		for i, waypoint := range pathfindingComponent.TargetCopy {
			if i > 0 {
				size := float32(4)
				from := prev
				to := math.Vector{waypoint.X - size/2, waypoint.Y - size/2}

				edge := to.Sub(from)
				angle := math.Atan232(edge.Y, edge.X)

				s.SpriteBatch.Draw(
					s.emptyTexture,
					from.X+size/2, from.Y+size/2, edge.Len(), 2,
					rendering.WithRotation(angle),
					rendering.WithOrigin(0, 1),
					rendering.WithColor(rendering.Color{0, 1, 1, 1}),
				)
			}

			prev = waypoint
		}

		if len(pathfindingComponent.Portals) > 0 {
			for _, p := range pathfindingComponent.Portals {
				size := float32(4)
				from := p.Left
				to := p.Right

				edge := to.Sub(from)
				angle := math.Atan232(edge.Y, edge.X)

				if from.Equal(to) {
					panic("WTH")
				}

				s.SpriteBatch.Draw(
					s.emptyTexture,
					from.X+size, from.Y+size, size*2, size*2,
					rendering.WithOrigin(size, size),
					rendering.WithColor(rendering.Color{1, 0.85, 0.73, 1}),
				)

				s.SpriteBatch.Draw(
					s.emptyTexture,
					to.X+size, to.Y+size, size*2, size*2,
					rendering.WithOrigin(size, size),
					rendering.WithColor(rendering.Color{1, 1, 0, 1}),
				)

				s.SpriteBatch.Draw(
					s.emptyTexture,
					from.X+size/2, from.Y+size/2, edge.Len(), 2,
					rendering.WithRotation(angle),
					rendering.WithOrigin(0, 1),
					rendering.WithColor(rendering.Color{0, 0, 1, 1}),
				)
			}
		}

		if debug {
			if len(pathfindingComponent.Path) > 0 {
				prev := pathfindingComponent.Path[0]

				for _, waypoint := range pathfindingComponent.Path[1:] {
					size := float32(4)
					from := prev
					to := math.Vector{waypoint.X - size/2, waypoint.Y - size/2}

					edge := to.Sub(from)
					angle := math.Atan232(edge.Y, edge.X)

					s.SpriteBatch.Draw(
						s.emptyTexture,
						from.X+size/2, from.Y+size/2, edge.Len(), 2,
						rendering.WithRotation(angle),
						rendering.WithOrigin(0, 1),
						rendering.WithColor(rendering.Color{1, 1, 0, 1}),
					)

					prev = waypoint
				}
			}
		}
	}

	s.SpriteBatch.End()
}
