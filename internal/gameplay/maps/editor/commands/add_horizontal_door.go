package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type AddHorizontalDoorMapCommand struct {
	m        *maps.TileMap
	row, col int
	backups  [2]int64
}

func NewAddHorizontalDoorMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewAddHorizontalDoorMapCommand, func(m *maps.TileMap, row, col int) bool {
		return m.GetAllBits(row, col, maps.FLOOR_BIT, maps.INTERIOR_WALL_S_BIT) && m.GetAllBits(row+1, col, maps.FLOOR_BIT, maps.INTERIOR_WALL_N_BIT)
	})
}

func NewAddHorizontalDoorMapCommand(m *maps.TileMap, row, col int) MapCommand {
	return &AddHorizontalDoorMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *AddHorizontalDoorMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row+1, c.col)

	c.m.SetBit(c.row, c.col, maps.DOOR_S_BIT)
	c.m.ClearBit(c.row, c.col, maps.INTERIOR_WALL_S_BIT)

	c.m.SetBit(c.row+1, c.col, maps.DOOR_N_BIT)
	c.m.ClearBit(c.row+1, c.col, maps.INTERIOR_WALL_N_BIT)
}

func (c *AddHorizontalDoorMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row+1, c.col, c.backups[1])
}
