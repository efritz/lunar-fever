package gameplay

type AddHorizontalWallMapCommand struct {
	m        *TileMap
	row, col int
	backups  [2]int64
}

func NewAddHorizontalWallMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewAddHorizontalWallMapCommand, func(m *TileMap, row, col int) bool {
		return m.GetBit(row, col, FLOOR_BIT) && !m.GetBit(row, col, INTERIOR_WALL_S_BIT) && m.GetBit(row+1, col, FLOOR_BIT) && !m.GetBit(row+1, col, INTERIOR_WALL_N_BIT)
	})
}

func NewAddHorizontalWallMapCommand(m *TileMap, row, col int) MapCommand {
	return &AddHorizontalWallMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *AddHorizontalWallMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row+1, c.col)

	c.m.SetBit(c.row, c.col, INTERIOR_WALL_S_BIT)
	c.m.ClearBit(c.row, c.col, DOOR_S_BIT)

	c.m.SetBit(c.row+1, c.col, INTERIOR_WALL_N_BIT)
	c.m.ClearBit(c.row+1, c.col, DOOR_N_BIT)
}

func (c *AddHorizontalWallMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row+1, c.col, c.backups[1])
}
