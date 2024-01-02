package gameplay

type RemoveVerticalWallMapCommand struct {
	m        *TileMap
	row, col int
	backups  [2]int64
}

func NewRemoveVerticalWallMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewRemoveVerticalWallMapCommand, func(m *TileMap, row, col int) bool {
		return m.GetAllBits(row, col, FLOOR_BIT, INTERIOR_WALL_E_BIT) && m.GetAllBits(row, col+1, FLOOR_BIT, INTERIOR_WALL_W_BIT)
	})
}

func NewRemoveVerticalWallMapCommand(m *TileMap, row, col int) MapCommand {
	return &RemoveVerticalWallMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *RemoveVerticalWallMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row, c.col+1)

	c.m.ClearBit(c.row, c.col, INTERIOR_WALL_E_BIT)
	c.m.ClearBit(c.row, c.col+1, INTERIOR_WALL_W_BIT)
}

func (c *RemoveVerticalWallMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row, c.col+1, c.backups[1])
}
