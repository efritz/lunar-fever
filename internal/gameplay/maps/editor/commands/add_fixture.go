package commands

import (
	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

type AddFixtureMapCommand struct {
	m        *maps.TileMap
	row, col int
}

func NewAddFixtureMapCommandFactory(fixture maps.Fixture) MapCommandFactory {
	return NewMapCommandFactory(NewAddFixtureMapCommand, func(m *maps.TileMap, row, col int) []TileIndex {
		set := pointsForFixture(row, col, fixture)

		for tileIndex := range set {
			// TODO - ooh also no floor bits in between
			if !m.GetBit(tileIndex.Row, tileIndex.Col, maps.FLOOR_BIT) {
				return nil
			}
		}

		for rowOffset := -fixtureExtents; rowOffset <= fixtureExtents; rowOffset++ {
			for colOffset := -fixtureExtents; colOffset <= fixture.TileHeight; colOffset++ {
				if row+rowOffset < 0 || row+rowOffset >= m.Height() || col+colOffset < 0 || col+colOffset >= m.Width() {
					continue
				}

				if fixture, ok := m.GetFixture(row+rowOffset, col+colOffset); ok {
					for tileIndex := range pointsForFixture(row+rowOffset, col+colOffset, fixture) {
						if _, ok := set[tileIndex]; ok {
							return nil
						}
					}
				}
			}
		}

		var tileIndexes []TileIndex
		for p := range set {
			tileIndexes = append(tileIndexes, p)
		}

		return tileIndexes
	})
}

func NewAddFixtureMapCommand(m *maps.TileMap, row, col int) MapCommand {
	return &AddFixtureMapCommand{
		m:   m,
		row: row,
		col: col,
	}
}

func (c *AddFixtureMapCommand) Execute() {
	fixture := maps.Fixtures[maps.FIXTURE_BENCH]
	c.m.SetFixture(c.row, c.col, fixture)

	for fixtureRow := 0; fixtureRow < fixture.TileHeight; fixtureRow++ {
		for fixtureCol := 0; fixtureCol < fixture.TileWidth; fixtureCol++ {
			row := c.row + fixtureRow
			col := c.col + fixtureCol

			if row-1 >= 0 && c.m.GetBit(row-1, col, maps.FLOOR_BIT) {
				c.m.SetBit(row-1, col, maps.FIXTURE_WALL_S_BIT)
				c.m.SetBit(row, col, maps.FIXTURE_WALL_N_BIT)
			}

			if row+1 < c.m.Height() && c.m.GetBit(row+1, col, maps.FLOOR_BIT) {
				c.m.SetBit(row+1, col, maps.FIXTURE_WALL_N_BIT)
				c.m.SetBit(row, col, maps.FIXTURE_WALL_S_BIT)
			}

			if col-1 >= 0 && c.m.GetBit(row, col-1, maps.FLOOR_BIT) {
				c.m.SetBit(row, col-1, maps.FIXTURE_WALL_E_BIT)
				c.m.SetBit(row, col, maps.FIXTURE_WALL_W_BIT)
			}

			if col+1 < c.m.Width() && c.m.GetBit(row, col+1, maps.FLOOR_BIT) {
				c.m.SetBit(row, col+1, maps.FIXTURE_WALL_W_BIT)
				c.m.SetBit(row, col, maps.FIXTURE_WALL_E_BIT)
			}
		}
	}
}

func (c *AddFixtureMapCommand) Unexecute() {
	c.m.SetFixture(c.row, c.col, maps.Fixtures[maps.FIXTURE_NONE])
}

const fixtureExtents = 8

func pointsForFixture(row, col int, fixture maps.Fixture) map[TileIndex]any {
	points := map[TileIndex]any{}
	for rowOffset := 0; rowOffset < fixture.TileHeight; rowOffset++ {
		for colOffset := 0; colOffset < fixture.TileWidth; colOffset++ {
			points[TileIndex{row + rowOffset, col + colOffset}] = struct{}{}
		}
	}

	return points
}
