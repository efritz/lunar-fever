package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type MapCommand interface {
	Execute()
	Unexecute()
}

type MapCommandFactory interface {
	Create(m *maps.TileMap, row, col int) MapCommand
	AffectedTileIndexes(m *maps.TileMap, row, col int) []TileIndex
}

type TileIndex struct {
	Row int
	Col int
}

type mapCommandFactory struct {
	createFunc              createFuncType
	affectedTileIndexesFunc affectedTileIndexesFuncType
}

type createFuncType func(m *maps.TileMap, row, col int) MapCommand
type affectedTileIndexesFuncType func(m *maps.TileMap, row, col int) []TileIndex

func NewMapCommandFactory(createFunc createFuncType, affectedTileIndexesFunc affectedTileIndexesFuncType) MapCommandFactory {
	return &mapCommandFactory{
		createFunc:              createFunc,
		affectedTileIndexesFunc: affectedTileIndexesFunc,
	}
}

func (f *mapCommandFactory) Create(m *maps.TileMap, row, col int) MapCommand {
	return f.createFunc(m, row, col)
}

func (f *mapCommandFactory) AffectedTileIndexes(m *maps.TileMap, row, col int) []TileIndex {
	return f.affectedTileIndexesFunc(m, row, col)
}
