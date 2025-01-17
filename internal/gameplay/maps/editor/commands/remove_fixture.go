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
}

func (c *RemoveFixtureMapCommand) Unexecute() {
	c.m.SetFixture(c.row, c.col, c.backup)
}
