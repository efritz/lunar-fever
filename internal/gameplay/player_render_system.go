package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type playerRenderSystem struct {
	*engine.Context
	playerCollection        *entity.Collection
	physicsComponentManager *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	emptyTexture            rendering.Texture
	headAtlas               rendering.Texture
	runAtlases              []rendering.Texture
}

func (s *playerRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
	s.headAtlas = s.TextureLoader.Load("character/headspack/head_1/head_1_1").Region(0, 0, 64, 64)
	s.runAtlases = []rendering.Texture{
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_1").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_2").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_3").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_4").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_5").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_6").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_7").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_8").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_9").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_10").Region(0, 0, 64, 64),
	}
}

func (s *playerRenderSystem) Exit() {}

var animationDelta int64
var animationIndex int64

func (s *playerRenderSystem) Process(elapsedMs int64) {
	animationDelta += elapsedMs
	for animationDelta > 150 {
		animationDelta -= 150
		animationIndex += 1
	}

	s.SpriteBatch.Begin()

	for _, entity := range s.playerCollection.Entities() {
		component, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		x1, y1, x2, y2 := component.Body.CoverBound()
		w := x2 - x1
		h := y2 - y1
		s.SpriteBatch.Draw(s.emptyTexture, x1, y1, w, h, rendering.WithColor(rendering.Color{1, 0, 1, .35}))

		x1, y1, x2, y2 = component.Body.NonorientedBound()
		w = x2 - x1
		h = y2 - y1
		s.SpriteBatch.Draw(s.emptyTexture, x1, y1, w, h, rendering.WithColor(rendering.Color{1, 1, 0, .35}), rendering.WithRotation(component.Body.Orient), rendering.WithOrigin(w/2, h/2))
		d := -float32(stdmath.Pi / 2)
		s.SpriteBatch.Draw(s.runAtlases[int(animationIndex)%len(s.runAtlases)], x1, y1, w, h, rendering.WithRotation(component.Body.Orient+d), rendering.WithOrigin(w/2, h/2))
		s.SpriteBatch.Draw(s.headAtlas, x1, y1, w, h, rendering.WithRotation(component.Body.Orient-d), rendering.WithOrigin(w/2, h/2))
	}

	s.SpriteBatch.End()
}
