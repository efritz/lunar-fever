package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type interactionRenderSystem struct {
	*GameContext
	emptyTexture rendering.Texture
}

func NewInteractionRenderSystem(ctx *GameContext) *interactionRenderSystem {
	return &interactionRenderSystem{GameContext: ctx}
}

func (s *interactionRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
}

func (s *interactionRenderSystem) Exit() {}

func (s *interactionRenderSystem) Process(elapsedMs int64) {
	for _, entity := range s.PlayerCollection.Entities() {
		physicsComponent, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		interactionComponent, ok := s.InteractionComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		healthComponent, ok := s.HealthComponentManager.GetComponent(entity)
		if !ok {
			return
		}

		if !interactionComponent.Interacting && canInteract(physicsComponent, interactionComponent, healthComponent) {
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
