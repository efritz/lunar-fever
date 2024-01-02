package gameplay

type RemoveFloorMapCommand struct {
	m        *TileMap
	row, col int
	backups  [5]int64
}

func NewRemoveFloorMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewRemoveFloorMapCommand, func(m *TileMap, row, col int) bool {
		return m.GetBit(row, col, FLOOR_BIT)
	})
}

func NewRemoveFloorMapCommand(m *TileMap, row, col int) MapCommand {
	return &RemoveFloorMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *RemoveFloorMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row-1, c.col)
	c.backups[2] = c.m.GetBits(c.row+1, c.col)
	c.backups[3] = c.m.GetBits(c.row, c.col-1)
	c.backups[4] = c.m.GetBits(c.row, c.col+1)

	if c.m.GetBit(c.row-1, c.col, FLOOR_BIT) {
		c.m.SetBit(c.row-1, c.col, INTERIOR_WALL_S_BIT)
	}

	if c.m.GetBit(c.row+1, c.col, FLOOR_BIT) {
		c.m.SetBit(c.row+1, c.col, INTERIOR_WALL_N_BIT)
	}

	if c.m.GetBit(c.row, c.col-1, FLOOR_BIT) {
		c.m.SetBit(c.row, c.col-1, INTERIOR_WALL_E_BIT)
	}

	if c.m.GetBit(c.row, c.col+1, FLOOR_BIT) {
		c.m.SetBit(c.row, c.col+1, INTERIOR_WALL_W_BIT)
	}

	c.m.ClearBit(c.row, c.col, FLOOR_BIT)
	c.m.ClearBit(c.row, c.col, INTERIOR_WALL_N_BIT)
	c.m.ClearBit(c.row, c.col, INTERIOR_WALL_S_BIT)
	c.m.ClearBit(c.row, c.col, INTERIOR_WALL_E_BIT)
	c.m.ClearBit(c.row, c.col, INTERIOR_WALL_W_BIT)
	c.m.ClearBit(c.row, c.col, DOOR_N_BIT)
	c.m.ClearBit(c.row, c.col, DOOR_S_BIT)
	c.m.ClearBit(c.row, c.col, DOOR_E_BIT)
	c.m.ClearBit(c.row, c.col, DOOR_W_BIT)

	c.m.ClearBit(c.row+1, c.col, DOOR_N_BIT)
	c.m.ClearBit(c.row-1, c.col, DOOR_S_BIT)
	c.m.ClearBit(c.row, c.col-1, DOOR_E_BIT)
	c.m.ClearBit(c.row, c.col+1, DOOR_W_BIT)
}

func (c *RemoveFloorMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row-1, c.col, c.backups[1])
	c.m.SetBits(c.row+1, c.col, c.backups[2])
	c.m.SetBits(c.row, c.col-1, c.backups[3])
	c.m.SetBits(c.row, c.col+1, c.backups[4])
}
