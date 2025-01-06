package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
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

func (s *npcMovementRenderSystem) Process(elapsedMs int64) {
	s.SpriteBatch.Begin()

	for _, entity := range s.NpcCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		pathfindingComponent, ok := s.PathfindingComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		prev := physicsComponent.Body.Position
		for _, waypoint := range pathfindingComponent.Target {
			size := float32(8)
			from := prev
			to := math.Vector{waypoint.X - size/2, waypoint.Y - size/2}

			var collisions []maps.Edge
			for _, obstacle := range s.Base.NavigationGraph.Obstacles {
				if waypoint.Equal(obstacle.From) || waypoint.Equal(obstacle.To) {
					collisions = append(collisions, obstacle)
				}
			}

			for _, co := range collisions {
				minX := math.Min(co.From.X, co.To.X)
				maxX := math.Max(co.From.X, co.To.X)
				minY := math.Min(co.From.Y, co.To.Y)
				maxY := math.Max(co.From.Y, co.To.Y)

				s.SpriteBatch.Draw(
					s.emptyTexture,
					minX, minY, maxX-minX+5, maxY-minY+5,
					rendering.WithColor(rendering.Color{1, 0, 0, 1}),
				)
			}

			edge := to.Sub(from)
			angle := math.Atan232(edge.Y, edge.X)

			s.SpriteBatch.Draw(
				s.emptyTexture,
				to.X+size/2, to.Y+size/2, size, size,
				rendering.WithOrigin(size/2, size/2),
				rendering.WithColor(rendering.Color{0, 1, 1, 1}),
			)

			s.SpriteBatch.Draw(
				s.emptyTexture,
				from.X+size/2, from.Y+size/2, edge.Len(), 2,
				rendering.WithRotation(angle),
				rendering.WithOrigin(0, 1),
				rendering.WithColor(rendering.Color{0, 1, 1, 1}),
			)

			prev = waypoint
		}
	}

	s.SpriteBatch.End()
}
