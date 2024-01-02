package system

import (
	"sort"

	"golang.org/x/exp/maps"
)

type Manager struct {
	initialized bool
	layers      []int
	systems     map[int][]System
}

func NewManager() *Manager {
	return &Manager{
		systems: map[int][]System{},
	}
}

func (m *Manager) Add(system System, layer int) {
	if m.initialized {
		system.Init()
	}

	if _, ok := m.systems[layer]; !ok {
		m.layers = append(m.layers, layer)
		sort.Ints(m.layers)
	}

	m.systems[layer] = append(m.systems[layer], system)
}

func (m *Manager) Init() {
	for _, layer := range m.layers {
		for _, s := range m.systems[layer] {
			s.Init()
		}
	}

	m.initialized = true
}

func (m *Manager) Exit() {
	for _, layer := range m.layers {
		for _, s := range m.systems[layer] {
			s.Exit()
		}
	}

	maps.Clear(m.systems)
	m.layers = nil
}

func (m *Manager) Process(elapsedMs int64) {
	for _, layer := range m.layers {
		for _, s := range m.systems[layer] {
			s.Process(elapsedMs)
		}
	}
}
