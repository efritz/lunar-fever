package commands

import "github.com/efritz/lunar-fever/internal/gameplay/maps"

type RemoveFixtureMapCommand struct {
	m        *maps.TileMap
	row, col int
	backup   maps.Fixture
}

func NewRemoveFixtureMapCommandFactory(fixture maps.Fixture) MapCommandFactory {
	return NewMapCommandFactory(NewRemoveFixtureMapCommand, func(m *maps.TileMap, row, col int) []TileIndex {
		if fixture, ok := m.GetFixture(row, col); ok {
			var tileIndexes []TileIndex
			for tileIndex := range pointsForFixture(row, col, fixture) {
				tileIndexes = append(tileIndexes, tileIndex)
			}

			return tileIndexes
		}

		return nil
	})
}

func NewRemoveFixtureMapCommand(m *maps.TileMap, row, col int) MapCommand {
	return &RemoveFixtureMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *RemoveFixtureMapCommand) Execute() {
	c.backup, _ = c.m.GetFixture(c.row, c.col)
	c.m.SetFixture(c.row, c.col, maps.Fixtures[maps.FIXTURE_NONE])

	if c.backup.TileWidth > 0 && c.backup.TileHeight > 0 {
		for fixtureRow := 0; fixtureRow < c.backup.TileHeight; fixtureRow++ {
			for fixtureCol := 0; fixtureCol < c.backup.TileWidth; fixtureCol++ {
				row := c.row + fixtureRow
				col := c.col + fixtureCol

				c.m.ClearBit(row, col, maps.FIXTURE_WALL_N_BIT)
				c.m.ClearBit(row, col, maps.FIXTURE_WALL_S_BIT)
				c.m.ClearBit(row, col, maps.FIXTURE_WALL_E_BIT)
				c.m.ClearBit(row, col, maps.FIXTURE_WALL_W_BIT)

				if row-1 >= 0 {
					c.m.ClearBit(row-1, col, maps.FIXTURE_WALL_S_BIT)
				}

				if row+1 < c.m.Height() {
					c.m.ClearBit(row+1, col, maps.FIXTURE_WALL_N_BIT)
				}

				if col-1 >= 0 {
					c.m.ClearBit(row, col-1, maps.FIXTURE_WALL_E_BIT)
				}

				if col+1 < c.m.Width() {
					c.m.ClearBit(row, col+1, maps.FIXTURE_WALL_W_BIT)
				}
			}
		}
	}
}

func (c *RemoveFixtureMapCommand) Unexecute() {
	c.m.SetFixture(c.row, c.col, c.backup)
}
