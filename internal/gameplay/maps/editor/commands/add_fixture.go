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
	c.m.SetFixture(c.row, c.col, maps.Fixtures[maps.FIXTURE_BENCH])
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
