package main

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/efritz/lunar-fever/internal/gameplay"
)

func main() {
	opts := []engine.DelegateOption{
		engine.WithInitialView(func(engineCtx *engine.Context) view.View {
			return view.NewTransitionView(gameplay.NewMainMenu(engineCtx), engineCtx.ViewManager)
		}),
	}

	game, err := engine.InitGame(opts...)
	if err != nil {
		panic(err)
	}

	game.Start()
}
