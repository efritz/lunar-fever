package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type npcRenderSystem struct {
	*engine.Context
	npcCollection           *entity.Collection
	physicsComponentManager *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	texture                 rendering.Texture
}

func (s *npcRenderSystem) Init() {
	s.texture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (s *npcRenderSystem) Exit() {}

func (s *npcRenderSystem) Process(elapsedMs int64) {
	s.SpriteBatch.Begin()

	for _, entity := range s.npcCollection.Entities() {
		component, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}
		body := component.Body

		x1, y1, x2, y2 := body.CoverBound()
		w := x2 - x1
		h := y2 - y1

		s.SpriteBatch.Draw(s.texture, x1, y1, w, h, rendering.WithColor(rendering.Color{1, 1, 0, 0.5}))

		for _, fixture := range body.Fixtures {
			for i := range fixture.Vertices {
				j := (i + 1) % len(fixture.Vertices)

				start := fixture.VertexInWorldSpace(body, i)
				end := fixture.VertexInWorldSpace(body, j)

				s.SpriteBatch.Draw(s.texture, start.X-3, start.Y-3, 6, 6, rendering.WithOrigin(5, 5))

				edge := end.Sub(start)
				angle := math.Atan232(edge.Y, edge.X)

				s.SpriteBatch.Draw(s.texture, start.X, start.Y, edge.Len(), 1, rendering.WithRotation(angle), rendering.WithOrigin(0, 1))
			}
		}

		// for _, fixture := range body.Fixtures {
		// 	for i := range fixture.Vertices {
		// 		j := (i + 1) % len(fixture.Vertices)

		// 		start := fixture.VertexInWorldSpace(body, i)
		// 		end := fixture.VertexInWorldSpace(body, j)

		// 		start = start.Add(end).Divs(2)
		// 		end = start.Add(fixture.NormalInWorldSpace(body, i).Muls(10))

		// 		// TODO - something is off here
		// 		edge := end.Sub(start)
		// 		angle := math.Atan232(edge.Y, edge.X)

		// 		s.SpriteBatch.Draw(s.texture, start.X, start.Y, edge.Len(), 1, rendering.WithColor(rendering.Color{1, 0, 0, 1}), rendering.WithRotation(angle), rendering.WithOrigin(0, 1))
		// 	}
		// }

		for _, fixture := range body.Fixtures {
			for i := range fixture.Vertices {
				j := (i + 1) % len(fixture.Vertices)

				start := fixture.VertexInWorldSpace(body, i)
				end := fixture.VertexInWorldSpace(body, j)

				start = start.Add(end).Divs(2)
				end = start.Add(end.Sub(start).Normalize().Muls(10))

				edge := end.Sub(start)
				angle := math.Atan232(edge.Y, edge.X)

				s.SpriteBatch.Draw(s.texture, start.X, start.Y, edge.Len(), 1, rendering.WithColor(rendering.Color{0, 0, 0, 1}), rendering.WithRotation(angle), rendering.WithOrigin(0, 1))
			}
		}

		for _, fixture := range body.Fixtures {
			for i := range fixture.Vertices {
				j := (i + 1) % len(fixture.Vertices)

				start := fixture.VertexInWorldSpace(body, i)
				end := fixture.VertexInWorldSpace(body, j)

				start = start.Add(end).Divs(2)
				end = start.Add(end.Sub(start).Normalize().Neg().Muls(10))

				edge := end.Sub(start)
				angle := math.Atan232(edge.Y, edge.X)

				s.SpriteBatch.Draw(s.texture, start.X, start.Y, edge.Len(), 1, rendering.WithColor(rendering.Color{0, 0, 0, 1}), rendering.WithRotation(angle), rendering.WithOrigin(1, 0))
			}
		}

		c2 := body.Position
		s.SpriteBatch.Draw(s.texture, c2.X-3, c2.Y-3, 6, 6, rendering.WithColor(rendering.Color{0, 0, 1, 1}))
	}

	s.SpriteBatch.End()
}
