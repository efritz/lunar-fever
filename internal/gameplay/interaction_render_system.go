package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type interactionRenderSystem struct {
	*engine.Context
	playerCollection            *entity.Collection
	physicsComponentManager     *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	interactionComponentManager *component.TypedManager[*InteractionComponent, InteractionComponentType]
	emptyTexture                rendering.Texture
}

func (s *interactionRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (s *interactionRenderSystem) Exit() {}

func (s *interactionRenderSystem) Process(elapsedMs int64) {
	for _, entity := range s.playerCollection.Entities() {
		physicsComponent, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		interactionComponent, ok := s.interactionComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		if !interactionComponent.Interacting && canInteract(physicsComponent, interactionComponent) {
			s.SpriteBatch.Begin()

			x1, y1, _, _ := physicsComponent.Body.NonorientedBound()
			w := float32(15)
			h := float32(15)

			s.SpriteBatch.Draw(s.emptyTexture, x1, y1, w, h, rendering.WithOrigin(w/2, h/2))
			s.SpriteBatch.End()

			font.Printf(x1+4, y1+12, "e",
				rendering.WithTextScale(0.3),
				rendering.WithTextColor(rendering.Color{0, 0, 0, 0.5}),
			)
		}
	}
}
