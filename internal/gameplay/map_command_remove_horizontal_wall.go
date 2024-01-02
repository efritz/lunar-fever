package gameplay

type RemoveHorizontalWallMapCommand struct {
	m        *TileMap
	row, col int
	backups  [2]int64
}

func NewRemoveHorizontalWallMapCommandFactory() MapCommandFactory {
	return NewMapCommandFactory(NewRemoveHorizontalWallMapCommand, func(m *TileMap, row, col int) bool {
		return m.GetAllBits(row, col, FLOOR_BIT, INTERIOR_WALL_S_BIT) && m.GetAllBits(row+1, col, FLOOR_BIT, INTERIOR_WALL_N_BIT)
	})
}

func NewRemoveHorizontalWallMapCommand(m *TileMap, row, col int) MapCommand {
	return &RemoveHorizontalWallMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *RemoveHorizontalWallMapCommand) Execute() {
	c.backups[0] = c.m.GetBits(c.row, c.col)
	c.backups[1] = c.m.GetBits(c.row+1, c.col)

	c.m.ClearBit(c.row, c.col, INTERIOR_WALL_S_BIT)
	c.m.ClearBit(c.row, c.col, DOOR_S_BIT)

	c.m.ClearBit(c.row+1, c.col, INTERIOR_WALL_N_BIT)
	c.m.ClearBit(c.row+1, c.col, DOOR_N_BIT)
}

func (c *RemoveHorizontalWallMapCommand) Unexecute() {
	c.m.SetBits(c.row, c.col, c.backups[0])
	c.m.SetBits(c.row+1, c.col, c.backups[1])
}
