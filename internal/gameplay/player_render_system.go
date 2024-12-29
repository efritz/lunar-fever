package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/ecs/component"
	"github.com/efritz/lunar-fever/internal/engine/ecs/entity"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type playerRenderSystem struct {
	*engine.Context
	playerCollection            *entity.Collection
	physicsComponentManager     *component.TypedManager[*physics.PhysicsComponent, physics.PhysicsComponentType]
	interactionComponentManager *component.TypedManager[*InteractionComponent, InteractionComponentType]
	healthComponentManager      *component.TypedManager[*HealthComponent, HealthComponentType]
	emptyTexture                rendering.Texture
	headAtlas                   rendering.Texture
	walkAtlases                 []rendering.Texture
	runAtlases                  []rendering.Texture
	interactAtlases             []rendering.Texture
	idleAtlas                   []rendering.Texture
	deathAtlas                  []rendering.Texture
	lastAnimationFrame          rendering.Texture
	animationQueue              *animationQueue
	distanceTraveled            float32
	wasMoving                   bool
	died                        bool

	idleTimer   int64
	idleRepeats int
}

func (s *playerRenderSystem) Init() {
	s.emptyTexture = s.TextureLoader.Load("base").Region(7*32, 1*32, 32, 32)
	s.headAtlas = s.TextureLoader.Load("character/headspack/head_1/head_1_1").Region(0, 0, 64, 64)
	s.walkAtlases = []rendering.Texture{
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_1").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_2").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_3").Region(0, 0, 64, 64), // idle
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_4").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_5").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_6").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_7").Region(0, 0, 64, 64), // idle
		s.TextureLoader.Load("character/scientist_1/walk_1/sci_walk_1_8").Region(0, 0, 64, 64),
	}
	s.runAtlases = []rendering.Texture{
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_1").Region(0, 0, 64, 64),  // corresopnds to walk_1_1
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_2").Region(0, 0, 64, 64),  // corresponds to walk_1_2
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_3").Region(0, 0, 64, 64),  // corresponds to walk_1_3 *
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_4").Region(0, 0, 64, 64),  // corresponds to walk_1_3 *
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_5").Region(0, 0, 64, 64),  // corresponds to walk_1_4
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_6").Region(0, 0, 64, 64),  // corresponds to walk_1_5
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_7").Region(0, 0, 64, 64),  // corresponds to walk_1_6
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_8").Region(0, 0, 64, 64),  // corresponds to walk_1_7 *
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_9").Region(0, 0, 64, 64),  // corresponds to walk_1_7 *
		s.TextureLoader.Load("character/scientist_1/run_1/sci_run_1_10").Region(0, 0, 64, 64), // corresponds to walk_1_8
	}
	s.interactAtlases = []rendering.Texture{
		s.TextureLoader.Load("character/scientist_1/interact_1/sci_interact_1_1").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/interact_1/sci_interact_1_2").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/interact_1/sci_interact_1_3").Region(0, 0, 64, 64),
	}
	s.idleAtlas = []rendering.Texture{
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_1").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_2").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_3").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_4").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_5").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_6").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_7").Region(0, 0, 64, 64),
		s.TextureLoader.Load("character/scientist_1/idle_2/sci_idle_2_8").Region(0, 0, 64, 64),
	}
	s.deathAtlas = []rendering.Texture{
		s.TextureLoader.Load("character/scientist_1/die_1/sci_fall_1_1").Region(0, 0, 64, 128),
		s.TextureLoader.Load("character/scientist_1/die_1/sci_fall_1_2").Region(0, 0, 64, 128),
		s.TextureLoader.Load("character/scientist_1/die_1/sci_fall_1_3").Region(0, 0, 64, 128),
		s.TextureLoader.Load("character/scientist_1/die_1/sci_fall_1_4").Region(0, 0, 64, 128),
		s.TextureLoader.Load("character/scientist_1/die_1/sci_fall_1_5").Region(0, 0, 64, 128),
		s.TextureLoader.Load("character/scientist_1/die_1/sci_fall_1_6").Region(0, 0, 64, 128),
	}

	s.lastAnimationFrame = s.walkAtlases[2]
	s.animationQueue = &animationQueue{}
}

func (s *playerRenderSystem) Exit() {}

const (
	minStartAnimationSpeed    = 0.05 // Minimum speed to start animation
	minContinueAnimationSpeed = 0.1  // Minimum speed to continue animation
	runThreshold              = 0.35 // Speed threshold to switch to run animation
	frameInterval             = 0.02 // How much distance to travel before next frame
	transitionSpeed           = 0.1  // How fast to transition between frames when stopping
	timeUntilIdleAnimation    = 3
	idleAnimationRepeats      = 3
)

func (s *playerRenderSystem) Process(elapsedMs int64) {
	s.SpriteBatch.Begin()

	for _, entity := range s.playerCollection.Entities() {
		physicsComponent, ok := s.physicsComponentManager.GetComponent(entity)
		if !ok {
			continue
		}

		interacting := false
		interactionComponent, ok := s.interactionComponentManager.GetComponent(entity)
		if ok {
			interacting = interactionComponent.Interacting
		}

		dead := false
		healthComponent, ok := s.healthComponentManager.GetComponent(entity)
		if ok {
			dead = healthComponent.Health <= 0
		}

		if dead {
			rotate := func(v math.Vector, angle float32) math.Vector {
				sinA := math.Sin32(angle)
				cosA := math.Cos32(angle)
				return math.Vector{
					X: v.X*cosA - v.Y*sinA,
					Y: v.X*sinA + v.Y*cosA,
				}
			}

			if physicsComponent.Body.Name == "player" {
				old := physicsComponent.Body
				physicsComponent.Body = physics.NewBody("player-dead", []physics.Fixture{
					physics.NewBasicFixture(
						0, 0, 48/2, 48, // bounds
						0.3, 0.2, // material
						0, 0, // friction
					),
				})
				physicsComponent.Body.Position = old.Position.Add(rotate(math.Vector{0, -(48 / 2)}, old.Orient))
				physicsComponent.Body.Orient = old.Orient
				physicsComponent.Body.Rotation = old.Rotation
			}
		}

		x1, y1, x2, y2 := physicsComponent.Body.NonorientedBound()
		w := x2 - x1
		h := y2 - y1

		if dead {
			if !s.died && s.animationQueue.Empty() {
				s.died = true
				s.animationQueue.Load(s.deathAtlas)
			}

			// Attempt to pop the next frame from the animation queue
			if frame, ok := s.animationQueue.Texture(elapsedMs); ok {
				s.lastAnimationFrame = frame
			}

			s.SpriteBatch.Draw(s.lastAnimationFrame, x1, y1, w, h, rendering.WithRotation(physicsComponent.Body.Orient), rendering.WithOrigin(w/2, h/2))
		} else {
			// Draw body
			s.SpriteBatch.Draw(s.selectBodyTexture(physicsComponent, interacting, elapsedMs), x1, y1, w, h, rendering.WithRotation(physicsComponent.Body.Orient), rendering.WithOrigin(w/2, h/2))

			// Always draw head
			s.SpriteBatch.Draw(s.headAtlas, x1, y1, w, h, rendering.WithRotation(physicsComponent.Body.Orient), rendering.WithOrigin(w/2, h/2))
		}
	}

	s.SpriteBatch.End()
}

func (s *playerRenderSystem) selectBodyTexture(component *physics.PhysicsComponent, interacting bool, elapsedMs int64) rendering.Texture {
	if speed := component.Body.LinearVelocity.Len(); s.canWalk(speed) {
		// Animate movement
		s.idleTimer = 0
		s.idleRepeats = 0
		s.wasMoving = true
		s.animationQueue.Reset()
		s.lastAnimationFrame = s.selectMovingTexture(speed, elapsedMs)
	} else {
		if s.wasMoving {
			// Transition to idle animation
			s.idleTimer = 0
			s.idleRepeats = 0
			s.wasMoving = false
			s.animationQueue.Load(s.selectPathToIdleFrame())
		}

		if interacting {
			// Poke animation
			s.idleTimer = 0
			s.idleRepeats = 0
			s.animationQueue.Load(s.interactAtlases)
		}

		s.idleTimer += elapsedMs
		if s.animationQueue.Empty() && s.idleTimer > timeUntilIdleAnimation*1000 {
			if s.idleRepeats < idleAnimationRepeats {
				s.idleRepeats++
				s.animationQueue.Load(s.idleAtlas)
			} else {
				s.idleTimer = 0
				s.idleRepeats = 0
			}
		}

		if s.animationQueue.Empty() {
			s.lastAnimationFrame = s.walkAtlases[2]
		}

		// Attempt to pop the next frame from the animation queue
		if frame, ok := s.animationQueue.Texture(elapsedMs); ok {
			s.lastAnimationFrame = frame
		}
	}

	return s.lastAnimationFrame
}

func (s *playerRenderSystem) canWalk(speed float32) bool {
	// Continuing movement
	if s.wasMoving && speed >= minContinueAnimationSpeed {
		return true
	}

	// Starting movement
	if !s.wasMoving && speed >= minStartAnimationSpeed {
		// Do not animate if we're animating the transition to idle
		return s.animationQueue.Empty()
	}

	// Velocity decaying toward zero
	return false
}

func (s *playerRenderSystem) selectMovingTexture(speed float32, elapsedMs int64) rendering.Texture {
	s.distanceTraveled += speed * float32(elapsedMs) / 1000.0
	progress := int(s.distanceTraveled / frameInterval)

	if speed >= runThreshold {
		return s.runAtlases[progress%len(s.runAtlases)]
	}

	return s.walkAtlases[progress%len(s.walkAtlases)]
}

func (s *playerRenderSystem) selectPathToIdleFrame() (path []rendering.Texture) {
	switch int(s.distanceTraveled/frameInterval) % len(s.walkAtlases) {
	case 0:
		// [0 ->] 1 -> 2
		return s.walkAtlases[1:3]
	case 1, 3:
		// [1,3 -> ] 2
		return s.walkAtlases[2:3]

	case 4:
		// [4 ->] 5 -> 6
		return s.walkAtlases[5:7]
	case 5, 7:
		// [5,7 ->] 6
		return s.walkAtlases[6:7]
	}

	// Empty path; already at idle position
	return nil
}

//
//

type animationQueue struct {
	timer  float32
	frames []rendering.Texture
}

func (q *animationQueue) Reset() {
	q.timer = 0
	q.frames = nil
}

func (q *animationQueue) Load(animations []rendering.Texture) {
	q.timer = transitionSpeed
	q.frames = animations
}

func (q *animationQueue) Empty() bool {
	return len(q.frames) == 0
}

func (q *animationQueue) Texture(elapsedMs int64) (rendering.Texture, bool) {
	if q.Empty() {
		return rendering.Texture{}, false
	}

	q.timer += float32(elapsedMs) / 1000.0

	if q.timer < transitionSpeed {
		return rendering.Texture{}, false
	}

	q.timer -= transitionSpeed
	frame := q.frames[0]
	q.frames = q.frames[1:]
	return frame, true
}
