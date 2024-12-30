package main

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/efritz/lunar-fever/internal/gameplay"
	"github.com/efritz/lunar-fever/internal/gameplay/menu"
)

func main() {
	opts := []engine.DelegateOption{
		engine.WithInitialView(func(engineCtx *engine.Context) view.View {
			return view.NewTransitionView(menu.NewMainMenu(engineCtx, gameplay.NewGameplay), engineCtx.ViewManager)
		}),
	}

	game, err := engine.InitGame(opts...)
	if err != nil {
		panic(err)
	}

	game.Start()
}
