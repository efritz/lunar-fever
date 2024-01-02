package view

import "slices"

type Manager struct {
	initialized bool
	views       []View
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) NumViews() int {
	return len(m.views)
}

func (m *Manager) Add(view View) {
	if m.initialized {
		view.Init()
	}

	m.views = append(m.views, view)
}

func (m *Manager) Remove(view View) {
	filtered := m.views[:0]
	for _, v := range m.views {
		if !areViewsEqual(v, view) {
			filtered = append(filtered, v)
		}
	}

	m.views = filtered
}

func (m *Manager) Clear() {
	for _, v := range m.views {
		if tv, ok := v.(*TransitionView); ok {
			tv.BeginExiting()
		} else {
			m.Remove(v)
		}
	}
}

func (m *Manager) Init() {
	for _, v := range m.views {
		v.Init()
	}

	m.initialized = true
}

func (m *Manager) Exit() {
	for _, v := range m.views {
		v.Exit()
	}

	m.views = m.views[:0]
}

func (m *Manager) Update(elapsedMs int64) {
	hasFocus := true

	// Create a copy of the view list to iterate. This allows a view to be removed during
	// a view's update method without concurrent modification issues.
	viewCopy := make([]View, len(m.views))
	copy(viewCopy, m.views)

	// Iterate the registered views in reverse order so that the top most (visible) view is
	// updated first.

	for i := len(viewCopy) - 1; i >= 0; i-- {
		v := viewCopy[i]

		// Do not process this view if it has been removed by another view's update method
		// during this tick.
		if !slices.Contains(m.views, v) {
			continue
		}

		v.Update(elapsedMs, hasFocus)

		// If this view is not an overlay, or if we've already encountered a non-overlay view,
		// set hasFocus to false for the next view.
		hasFocus = hasFocus && v.IsOverlay()
	}
}

func (m *Manager) Render(elapsedMs int64) {
	for _, v := range m.views {
		v.Render(elapsedMs)
	}
}
