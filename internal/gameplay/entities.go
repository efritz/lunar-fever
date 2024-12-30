package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/physics"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

func createPlayer(ctx *GameContext) {
	player := ctx.EntityManager.Create()
	ctx.TagManager.SetTag(player, "player")
	ctx.GroupManager.AddGroup(player, "scientist")
	ctx.GroupManager.AddGroup(player, "physics")

	body := physics.NewBody("scientist", []physics.Fixture{
		physics.NewBasicFixture(
			0, 0, 48/2, 48/2, // bounds
			0.3, 0.2, // material
			0, 0, // friction
		),
	})
	body.Position = math.Vector{rendering.DisplayWidth - 200, 400}
	ctx.PhysicsComponentManager.AddComponent(player, &physics.PhysicsComponent{Body: body})
	ctx.InteractionComponentManager.AddComponent(player, &InteractionComponent{})
	ctx.HealthComponentManager.AddComponent(player, &HealthComponent{Health: 100, MaxHealth: 100})
}

func createScientist(ctx *GameContext) {
	player := ctx.EntityManager.Create()
	ctx.GroupManager.AddGroup(player, "scientist")
	ctx.GroupManager.AddGroup(player, "physics")

	body := physics.NewBody("scientist", []physics.Fixture{
		physics.NewBasicFixture(
			0, 0, 48/2, 48/2, // bounds
			0.3, 0.2, // material
			0, 0, // friction
		),
	})
	body.Position = math.Vector{rendering.DisplayWidth - 100, 300}
	ctx.PhysicsComponentManager.AddComponent(player, &physics.PhysicsComponent{Body: body})
	ctx.HealthComponentManager.AddComponent(player, &HealthComponent{Health: 100, MaxHealth: 100})
}

func createRover(ctx *GameContext) {
	rover := ctx.EntityManager.Create()
	ctx.TagManager.SetTag(rover, "rover")
	ctx.GroupManager.AddGroup(rover, "physics")

	body := physics.NewBody("rover", []physics.Fixture{
		physics.NewBasicFixture(
			0, 0, 69, 123, // bounds
			20, 0.5, // material
			0, 0, // friction
		),
	})
	body.Position = math.Vector{rendering.DisplayWidth / 4, rendering.DisplayHeight / 4}
	ctx.PhysicsComponentManager.AddComponent(rover, &physics.PhysicsComponent{Body: body})
}

func createNPC(ctx *GameContext) {
	npc := ctx.EntityManager.Create()
	ctx.GroupManager.AddGroup(npc, "npc")
	ctx.GroupManager.AddGroup(npc, "physics")

	points1 := make([]math.Vector, 5)
	points2 := make([]math.Vector, 5)
	points3 := make([]math.Vector, 5)

	tx := float32(40)
	ty := float32(40)
	hw := float32(32)
	hh := float32(48)

	for _, vectors := range [][]math.Vector{points1, points2, points3} {
		tx += hw * 2.25
		ty += hh / 2

		vectors[0] = math.Vector{tx, ty + hh*2}
		vectors[1] = math.Vector{tx + hw, ty - hh}
		vectors[2] = math.Vector{tx - hw, ty + hh}
		vectors[3] = math.Vector{tx + hw, ty + hh}
		vectors[4] = math.Vector{tx - hw, ty - hh}
	}

	body := physics.NewBody("npc", []physics.Fixture{
		physics.NewFixture(points1, 0.1, 0.1, 0, 0),
		physics.NewFixture(points2, 0.1, 0.1, 0, 0),
		physics.NewFixture(points3, 0.4, 0.1, 0, 0),
	})
	ctx.PhysicsComponentManager.AddComponent(npc, &physics.PhysicsComponent{Body: body})
}

func createWalls(ctx *GameContext) {
	type Options struct {
		name    string
		w       float32
		h       float32
		iOffset float32
		jOffset float32
	}

	build := func(i, j int, opts Options) {
		entity := ctx.EntityManager.Create()
		ctx.GroupManager.AddGroup(entity, "physics")
		ctx.GroupManager.AddGroup(entity, opts.name)

		body := physics.NewBody(opts.name, []physics.Fixture{
			physics.NewBasicFixture(
				0, 0, opts.w, opts.h, // bounds
				0.0, 0.5, // material
				0, 0, // friction
			),
		})
		body.Position = math.Vector{float32(j*64) + opts.jOffset, float32(i*64) + opts.iOffset}
		ctx.PhysicsComponentManager.AddComponent(entity, &physics.PhysicsComponent{Body: body})
	}

	parametersByBit := map[maps.TileBitIndex]Options{
		maps.INTERIOR_WALL_N_BIT: {"wall", 32, 2, +1, 32},
		maps.INTERIOR_WALL_S_BIT: {"wall", 32, 2, 64 - 1, 32},
		maps.INTERIOR_WALL_W_BIT: {"wall", 2, 32, 32, +1},
		maps.INTERIOR_WALL_E_BIT: {"wall", 2, 32, 32, 64 - 1},
		maps.DOOR_N_BIT:          {"door", 32, 2, +1, 32},
		maps.DOOR_S_BIT:          {"door", 32, 2, 64 - 1, 32},
		maps.DOOR_W_BIT:          {"door", 2, 32, 32, +1},
		maps.DOOR_E_BIT:          {"door", 2, 32, 32, 64 - 1},
	}

	for i := 0; i < ctx.TileMap.Height(); i++ {
		for j := 0; j < ctx.TileMap.Width(); j++ {
			for bit, opts := range parametersByBit {
				if ctx.TileMap.GetBit(i, j, bit) {
					build(i, j, opts)
				}
			}
		}
	}
}
