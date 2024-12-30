package menu

import (
	"github.com/efritz/lunar-fever/internal/engine"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type MenuDelegate interface {
	Init()
	Exit()
	Update(elapsedMs int64, hasFocus bool)
	Render(elapsedMs int64)
}

type Menu struct {
	*engine.Context
	delegate    MenuDelegate
	initialized bool
	selected    int
	entries     []*MenuEntry
}

func NewMenu(engineCtx *engine.Context, delegate MenuDelegate) *Menu {
	return &Menu{
		Context:  engineCtx,
		delegate: delegate,
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
	if m.delegate != nil {
		m.delegate.Init()
	}

	for _, entry := range m.entries {
		entry.Init()
	}

	m.initialized = true
	initFonts()
}

func (m *Menu) Exit() {
	if m.delegate != nil {
		m.delegate.Exit()
	}

	for _, entry := range m.entries {
		entry.Exit()
	}
}

func (m *Menu) Update(elapsedMs int64, hasFocus bool) {
	if m.delegate != nil {
		m.delegate.Update(elapsedMs, hasFocus)
	}

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
	if m.delegate != nil {
		m.delegate.Render(elapsedMs)
	}

	for i, entry := range m.entries {
		entry.Render(elapsedMs, i, i == m.selected)
	}
}

func (m *Menu) IsOverlay() bool {
	return false
}
