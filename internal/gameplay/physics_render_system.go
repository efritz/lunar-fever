package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type physicsRenderSystem struct {
	*engine.Context
	entityCollection        *entity.Collection
	physicsComponentManager *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	emptyTexture            rendering.Texture
}

func (s *physicsRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (s *physicsRenderSystem) Exit() {}

func (s *physicsRenderSystem) Process(elapsedMs int64) {
	if !debug {
		return
	}

	s.SpriteBatch.Begin()

	for _, entity := range s.entityCollection.Entities() {
		component, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok || component.CollisionsDisabled {
			continue
		}

		x01, y01, x02, y02 := component.Body.CoverBound()
		w := x02 - x01
		h := y02 - y01

		s.SpriteBatch.Draw(s.emptyTexture, x01, y01, w, h, rendering.WithColor(rendering.Color{1, 0, 1, .35}))

		x1, y1, x2, y2 := component.Body.NonorientedBound()
		w = x2 - x1
		h = y2 - y1
		s.SpriteBatch.Draw(s.emptyTexture, x1, y1, w, h, rendering.WithColor(rendering.Color{1, 1, 0, .35}), rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w/2, h/2))

	}

	s.SpriteBatch.End()
}
