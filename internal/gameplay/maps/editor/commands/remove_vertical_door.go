package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type RemoveVerticalDoorMapCommand struct {
	m        *maps.TileMap
	row, col int
	backups  [2]int64
}

func NewRemoveVerticalDoorMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewRemoveVerticalDoorMapCommand, func(m *maps.TileMap, row, col int) bool {
		return m.GetAllBits(row, col, maps.FLOOR_BIT, maps.DOOR_E_BIT) && m.GetAllBits(row, col+1, maps.FLOOR_BIT, maps.DOOR_W_BIT)
	})
}

func NewRemoveVerticalDoorMapCommand(m *maps.TileMap, row, col int) MapCommand {
	return &RemoveVerticalDoorMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *RemoveVerticalDoorMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row, c.col+1)

	c.m.ClearBit(c.row, c.col, maps.DOOR_E_BIT)
	c.m.SetBit(c.row, c.col, maps.INTERIOR_WALL_E_BIT)

	c.m.ClearBit(c.row, c.col+1, maps.DOOR_W_BIT)
	c.m.SetBit(c.row, c.col+1, maps.INTERIOR_WALL_W_BIT)
}

func (c *RemoveVerticalDoorMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row, c.col+1, c.backups[1])
}
