package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type RemoveVerticalWallMapCommand struct {
	m        *maps.TileMap
	row, col int
	backups  [2]int64
}

func NewRemoveVerticalWallMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewRemoveVerticalWallMapCommand, func(m *maps.TileMap, row, col int) []TileIndex {
		if m.GetAllBits(row, col, maps.FLOOR_BIT, maps.INTERIOR_WALL_E_BIT) && m.GetAllBits(row, col+1, maps.FLOOR_BIT, maps.INTERIOR_WALL_W_BIT) {
			return []TileIndex{{row, col}}
		}

		return nil
	})
}

func NewRemoveVerticalWallMapCommand(m *maps.TileMap, row, col int) MapCommand {
	return &RemoveVerticalWallMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *RemoveVerticalWallMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row, c.col+1)

	c.m.ClearBit(c.row, c.col, maps.INTERIOR_WALL_E_BIT)
	c.m.ClearBit(c.row, c.col+1, maps.INTERIOR_WALL_W_BIT)
}

func (c *RemoveVerticalWallMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row, c.col+1, c.backups[1])
}
