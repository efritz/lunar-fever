package engine

import (
	"github.com/efritz/lunar-fever/internal/engine/camera"
	"github.com/efritz/lunar-fever/internal/engine/game"
	"github.com/efritz/lunar-fever/internal/engine/input"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
	"github.com/efritz/lunar-fever/internal/engine/view"
)

type Context struct {
	RenderingContext
	InputContext
	Game        *game.Game
	ViewManager *view.Manager
	Camera      *camera.Camera // TODO - extract from engine completely?
}

type (
	RenderingContext = rendering.Context
	InputContext     = input.Context
)

func InitContext(game *game.Game) (*Context, error) {
	renderingCtx, err := rendering.InitContext()
	if err != nil {
		return nil, err
	}

	return &Context{
		RenderingContext: renderingCtx,
		InputContext:     input.InitContext(renderingCtx.Window),
		Game:             game,
		ViewManager:      view.NewManager(),
		Camera:           camera.NewCamera(),
	}, nil
}
