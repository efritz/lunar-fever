package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type MapCommand interface {
	Execute()
	Unexecute()
}

type MapCommandFactory interface {
	Create(m *maps.TileMap, row, col int) MapCommand
	Valid(m *maps.TileMap, row, col int) bool
}

type mapCommandFactory struct {
	createFunc createFuncType
	validFunc  validFuncType
}

type createFuncType func(m *maps.TileMap, row, col int) MapCommand
type validFuncType func(m *maps.TileMap, row, col int) bool

func NewMapCommandFactory(createFunc createFuncType, validFunc validFuncType) MapCommandFactory {
	return &mapCommandFactory{
		createFunc: createFunc,
		validFunc:  validFunc,
	}
}

func (f *mapCommandFactory) Create(m *maps.TileMap, row, col int) MapCommand {
	return f.createFunc(m, row, col)
}

func (f *mapCommandFactory) Valid(m *maps.TileMap, row, col int) bool {
	return f.validFunc(m, row, col)
}
