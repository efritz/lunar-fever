package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type regolithRenderSystem struct {
	*GameContext
	textures []rendering.Texture
}

func NewRegolithRenderSystem(ctx *GameContext) *regolithRenderSystem {
	return &regolithRenderSystem{GameContext: ctx}
}

func (s *regolithRenderSystem) Init() {
	texture := s.TextureLoader.Load("regolith")
	s.textures = append(s.textures, texture.Region(0, 0, 71, 71))
	s.textures = append(s.textures, texture.Region(0, 71, 141, 141))
	s.textures = append(s.textures, texture.Region(0, 71+141, 373, 373))
}

func (s *regolithRenderSystem) Exit() {}

func (s *regolithRenderSystem) Process(elapsedMs int64) {
	s.SpriteBatch.Begin()
	for _, texture := range s.textures {
		s.drawRegolith(texture)
	}
	s.SpriteBatch.End()
}

func (s *regolithRenderSystem) drawRegolith(texture rendering.Texture) {
	x1, y1, x2, y2 := s.Camera.Bounds()

	w := int((texture.U2 - texture.U1) * texture.Width)
	h := int((texture.V2 - texture.V1) * texture.Height)

	startX := math.PrevMultiple(int(x1), int(w))
	startY := math.PrevMultiple(int(y1), int(h))

	for i := 0; startX+i <= int(x2); i += w {
		for j := 0; startY+j <= int(y2); j += h {
			s.SpriteBatch.Draw(texture, float32(startX+i), float32(startY+j), float32(w), float32(h))
		}
	}
}
