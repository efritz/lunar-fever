package gameplay

import (
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Menu struct {
	*engine.Context
	initialized bool
	selected    int
	entries     []*MenuEntry
}

func NewMenu(engineCtx *engine.Context) *Menu {
	return &Menu{
		Context: engineCtx,
	}
}

func (m *Menu) AddEntry(text string, delegate MenuEntrySelectionDelegate) {
	menuEntry := NewMenuEntry(m.Context, text, delegate)
	m.entries = append(m.entries, menuEntry)

	if m.initialized {
		menuEntry.Init()
	}
}

func (m *Menu) Init() {
	for _, entry := range m.entries {
		entry.Init()
	}

	m.initialized = true
	initFonts()
}

func (m *Menu) Exit() {
	for _, entry := range m.entries {
		entry.Exit()
	}
}

func (m *Menu) Update(elapsedMs int64, hasFocus bool) {
	if m.Keyboard.IsKeyNewlyDown(glfw.KeyUp) {
		m.selected--
		if m.selected < 0 {
			m.selected = len(m.entries) - 1
		}
	}

	if m.Keyboard.IsKeyNewlyDown(glfw.KeyDown) {
		m.selected++
		if m.selected >= len(m.entries) {
			m.selected = 0
		}
	}

	if m.Keyboard.IsKeyNewlyDown(glfw.KeyEnter) {
		m.entries[m.selected].OnSelect()
	}

	for i, entry := range m.entries {
		entry.Update(elapsedMs, i == m.selected)
	}
}

func (m *Menu) Render(elapsedMs int64) {
	offset := int64(128)
	for _, entry := range m.entries {
		entry.SetPosition(64, offset)
		offset += 64
	}

	for i, entry := range m.entries {
		entry.Render(elapsedMs, i == m.selected)
	}
}

func (m *Menu) IsOverlay() bool {
	return false
}

var font *rendering.Font

func initFonts() {
	var err error
	if font, err = rendering.LoadFont("Roboto-Light"); err != nil {
		panic(err)
	}
}
