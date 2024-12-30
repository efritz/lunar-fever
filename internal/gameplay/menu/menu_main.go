package menu

import (
	"fmt"
	"time"

	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/view"
	"github.com/efritz/lunar-fever/internal/gameplay/maps/editor"
)

func NewMainMenu(engineCtx *engine.Context, gameplayFactory func(*engine.Context) view.View) view.View {
	updater, err := NewUpdater()
	if err != nil {
		panic(err)
	}

	menu := NewMenu(engineCtx, nil)
	hasUpdate, err := updater.HasUpdate()
	if err != nil {
		panic(err)
	}
	if hasUpdate {
		menu.AddEntry("Download update", &downloadUpdateMenuEntry{updater: updater})
	}
	menu.AddEntry("Play", &gameplayMenuEntry{Context: engineCtx, gameplayFactory: gameplayFactory})
	menu.AddEntry("Tile editor", &tileEditorMenuEntry{Context: engineCtx})
	menu.AddEntry("Exit", &exitMenuEntry{exit: engineCtx.Game.Stop})

	return menu
}

type downloadUpdateMenuEntry struct {
	updater *updater
}

func (e *downloadUpdateMenuEntry) OnSelect() {
	if err := e.updater.Update(); err != nil {
		panic(err)
	}
}

type gameplayMenuEntry struct {
	*engine.Context
	gameplayFactory func(*engine.Context) view.View
}

func (e *gameplayMenuEntry) OnSelect() {
	Load(e.Context, e.gameplayFactory(e.Context), fakeLoader)
}

type tileEditorMenuEntry struct {
	*engine.Context
}

func (e *tileEditorMenuEntry) OnSelect() {
	Load(e.Context, editor.NewEditor(e.Context), fakeLoader)
}

var fakeLoader = func() {
	total := 5
	interval := time.Second / 4 * 0

	for i := 0; i < total; i++ {
		fmt.Printf("Loading %d of %d...\n", i+1, total)
		time.Sleep(interval)
	}

	fmt.Printf("Done!\n")
}

type exitMenuEntry struct {
	exit func()
}

func (e *exitMenuEntry) OnSelect() {
	e.exit()
}
