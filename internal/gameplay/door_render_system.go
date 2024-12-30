package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/engine/ecs/system"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type doorRenderSystem struct {
	*GameContext
	texture rendering.Texture
}

func NewDoorRenderSystem(ctx *GameContext) system.System {
	return &doorRenderSystem{GameContext: ctx}
}

func (s *doorRenderSystem) Init() {
	s.texture = s.TextureLoader.Load("base").Region(3*64, 128-2, 64, 2)
}

func (s *doorRenderSystem) Exit() {}

func (s *doorRenderSystem) Process(elapsedMs int64) {
	s.SpriteBatch.Begin()

	for _, entity := range s.DoorCollection.Entities() {
		component, ok := s.PhysicsComponentManager.GetComponent(entity)
		if !ok || component.CollisionsDisabled {
			continue
		}

		x01, y01, x02, y02 := component.Body.CoverBound()
		w := x02 - x01
		h := y02 - y01

		opts := []rendering.DrawOptionFunc{}
		if h > w {
			w, h = h, w
			opts = append(opts,
				rendering.WithOrigin(32, 32),
				rendering.WithRotation((90+180)*stdmath.Pi/180),
			)
		}

		s.SpriteBatch.Draw(s.texture, x01, y01, w, h, opts...)
	}

	s.SpriteBatch.End()
}
