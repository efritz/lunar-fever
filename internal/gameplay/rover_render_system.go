package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type roverRenderSystem struct {
	*GameContext
	emptyTexture rendering.Texture
	baseTexture  rendering.Texture
	axleTexture  rendering.Texture
	tireTexture  rendering.Texture
}

func NewRoverRenderSystem(ctx *GameContext) system.System {
	return &roverRenderSystem{GameContext: ctx}
}

func (s *roverRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
	s.baseTexture = s.TextureLoader.Load("rover/base").Region(0, 0, 69, 123)
	s.axleTexture = s.TextureLoader.Load("rover/axle").Region(0, 0, 177, 45)
	s.tireTexture = s.TextureLoader.Load("rover/tire").Region(0, 0, 22, 45)
}

func (s *roverRenderSystem) Exit() {}

var (
	// TODO - deglobalize
	animationDelta int64
	animationIndex int64
)

func (s *roverRenderSystem) Process(elapsedMs int64) {
	animationDelta += elapsedMs
	for animationDelta > 150 {
		animationDelta -= 150
		animationIndex += 1
	}

	s.SpriteBatch.Begin()

	for _, entity := range s.RoverCollection.Entities() {
		component, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		x01, y01, x02, y02 := component.Body.CoverBound()
		w := x02 - x01
		h := y02 - y01

		// s.SpriteBatch.Draw(s.emptyTexture, x01, y01, w, h, rendering.WithColor(rendering.Color{1, 0, 1, .35}))

		x1, y1, x2, y2 := component.Body.NonorientedBound()
		w = x2 - x1
		h = y2 - y1
		// s.SpriteBatch.Draw(s.emptyTexture, x1, y1, w, h, rendering.WithColor(rendering.Color{1, 1, 0, .35}), rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w/2, h/2))

		w2 := float32(177)
		h2 := float32(45)
		axleOffset := float32(60)
		w3 := float32(22)
		h3 := h2

		v2 := component.Body.Rotation.Mul(math.Vector{-w2 / 2, -axleOffset}).Add(math.Vector{x1 + w/2 - w3, y1 + h/2 - h3/2})
		v3 := component.Body.Rotation.Mul(math.Vector{+w2 / 2, -axleOffset}).Add(math.Vector{x1 + w/2, y1 + h/2 - h3/2})

		// Axel and body
		s.SpriteBatch.Draw(s.axleTexture, (x1+w/2)-w2/2, (y1+h/2)-h2/2+axleOffset, w2, h2, rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w2/2, h2/2-axleOffset))
		s.SpriteBatch.Draw(s.axleTexture, (x1+w/2)-w2/2, (y1+h/2)-h2/2-axleOffset, w2, h2, rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w2/2, h2/2+axleOffset))
		s.SpriteBatch.Draw(s.baseTexture, x1, y1, w, h, rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w/2, h/2))

		// Front tires
		s.SpriteBatch.Draw(s.tireTexture, v2.X, v2.Y, w3, h3, rendering.WithRotation(component.Body.Orient+tireRotation), rendering.WithOrigin(w3, h3/2))
		s.SpriteBatch.Draw(s.tireTexture, v3.X, v3.Y, w3, h3, rendering.WithRotation(component.Body.Orient+tireRotation), rendering.WithOrigin(0, h3/2), rendering.WithSpriteEffects(rendering.SpriteEffectFlipHorizontal))

		// Back tires
		s.SpriteBatch.Draw(s.tireTexture, (x1+w/2)-w2/2-w3, (y1+h/2)-h3/2+axleOffset, w3, h3, rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w2/2+w3, h3/2-axleOffset))
		s.SpriteBatch.Draw(s.tireTexture, (x1+w/2)+w2/2, (y1+h/2)-h3/2+axleOffset, w3, h3, rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(-w2/2, h3/2-axleOffset), rendering.WithSpriteEffects(rendering.SpriteEffectFlipHorizontal))
	}

	s.SpriteBatch.End()
}

// TODO - deglobalize
var tireRotation float32
