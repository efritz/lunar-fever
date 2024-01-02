package gameplay

import (
	"fmt"
	"time"

	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/view"
)

func NewMainMenu(engineCtx *engine.Context) view.View {
	menu := NewMenu(engineCtx)
	menu.AddEntry(&gameplayMenuEntry{Context: engineCtx})
	menu.AddEntry(&tileEditorMenuEntry{Context: engineCtx})
	menu.AddEntry(&noopMenuEntry{})
	menu.AddEntry(&exitMenuEntry{exit: engineCtx.Game.Stop})

	return menu
}

type gameplayMenuEntry struct {
	*engine.Context
}

func (e *gameplayMenuEntry) OnSelect() {
	Load(e.Context, NewGameplay(e.Context), fakeLoader)
}

type tileEditorMenuEntry struct {
	*engine.Context
}

func (e *tileEditorMenuEntry) OnSelect() {
	Load(e.Context, NewEditor(e.Context), fakeLoader)
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

type noopMenuEntry struct{}

func (e *noopMenuEntry) OnSelect() {
	fmt.Printf("TEST!\n")
}

type exitMenuEntry struct {
	exit func()
}

func (e *exitMenuEntry) OnSelect() {
	e.exit()
}
