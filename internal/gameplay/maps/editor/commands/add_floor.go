package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type AddFloorMapCommand struct {
	m        *maps.TileMap
	row, col int
	backups  [5]int64
}

func NewAddFloorMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewAddFloorMapCommand, func(m *maps.TileMap, row, col int) []TileIndex {
		if !m.GetBit(row, col, maps.FLOOR_BIT) {
			return []TileIndex{{row, col}}
		}

		return nil
	})
}

func NewAddFloorMapCommand(m *maps.TileMap, row, col int) MapCommand {
	return &AddFloorMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *AddFloorMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row-1, c.col)
	c.backups[2] = c.m.GetBits(c.row+1, c.col)
	c.backups[3] = c.m.GetBits(c.row, c.col-1)
	c.backups[4] = c.m.GetBits(c.row, c.col+1)

	if c.m.GetBit(c.row-1, c.col, maps.FLOOR_BIT) {
		c.m.ClearBit(c.row-1, c.col, maps.INTERIOR_WALL_S_BIT)
	} else {
		c.m.SetBit(c.row, c.col, maps.INTERIOR_WALL_N_BIT)
	}

	if c.m.GetBit(c.row+1, c.col, maps.FLOOR_BIT) {
		c.m.ClearBit(c.row+1, c.col, maps.INTERIOR_WALL_N_BIT)
	} else {
		c.m.SetBit(c.row, c.col, maps.INTERIOR_WALL_S_BIT)
	}

	if c.m.GetBit(c.row, c.col-1, maps.FLOOR_BIT) {
		c.m.ClearBit(c.row, c.col-1, maps.INTERIOR_WALL_E_BIT)
	} else {
		c.m.SetBit(c.row, c.col, maps.INTERIOR_WALL_W_BIT)
	}

	if c.m.GetBit(c.row, c.col+1, maps.FLOOR_BIT) {
		c.m.ClearBit(c.row, c.col+1, maps.INTERIOR_WALL_W_BIT)
	} else {
		c.m.SetBit(c.row, c.col, maps.INTERIOR_WALL_E_BIT)
	}

	c.m.SetBit(c.row, c.col, maps.FLOOR_BIT)
}

func (c *AddFloorMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row-1, c.col, c.backups[1])
	c.m.SetBits(c.row+1, c.col, c.backups[2])
	c.m.SetBits(c.row, c.col-1, c.backups[3])
	c.m.SetBits(c.row, c.col+1, c.backups[4])
}
