package maps

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type BaseRenderer struct {
	spriteBatch   *rendering.SpriteBatch
	tileMap       *TileMap
	textures      map[TileBitIndex]baseTexture
	emptyTexture  rendering.Texture
	fixturesAtlas rendering.Texture
}

func NewBaseRenderer(spriteBatch *rendering.SpriteBatch, textureLoader *rendering.TextureLoader, tileMap *TileMap, renderDoors bool) *BaseRenderer {
	return &BaseRenderer{
		spriteBatch:   spriteBatch,
		tileMap:       tileMap,
		textures:      newBaseTextureMap(textureLoader.Load("base"), renderDoors),
		emptyTexture:  textureLoader.Load("base").Region(7*32, 1*32, 32, 32),
		fixturesAtlas: textureLoader.Load("fixtures"),
	}
}

func setAestheticBits(tileMap *TileMap) *TileMap {
	m2 := NewTileMap(tileMap.Width(), tileMap.Height(), tileMap.GridSize())
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			m2.SetBits(row, col, tileMap.GetBits(row, col))
		}
	}

	setExternalWallTiles(m2)
	setCornerTiles(m2)
	return m2
}

func setExternalWallTiles(tileMap *TileMap) {
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row-1, col, FLOOR_BIT) {
				tileMap.SetBit(row-1, col, EXTERIOR_WALL_S_BIT)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row+1, col, FLOOR_BIT) {
				tileMap.SetBit(row+1, col, EXTERIOR_WALL_N_BIT)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
				tileMap.SetBit(row, col+1, EXTERIOR_WALL_W_BIT)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) {
				tileMap.SetBit(row, col-1, EXTERIOR_WALL_E_BIT)
			}
		}
	}
}

func setCornerTiles(tileMap *TileMap) {
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if tileMap.GetBit(row-1, col, EXTERIOR_WALL_S_BIT) && tileMap.GetBit(row, col+1, EXTERIOR_WALL_W_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
				tileMap.SetBit(row-1, col+1, EXTERIOR_CORNER_CONCAVE_SW_BIT)
			}

			if tileMap.GetBit(row-1, col, EXTERIOR_WALL_S_BIT) && tileMap.GetBit(row, col-1, EXTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) {
				tileMap.SetBit(row-1, col-1, EXTERIOR_CORNER_CONCAVE_SE_BIT)
			}

			if tileMap.GetBit(row+1, col, EXTERIOR_WALL_N_BIT) && tileMap.GetBit(row, col+1, EXTERIOR_WALL_W_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
				tileMap.SetBit(row+1, col+1, EXTERIOR_CORNER_CONCAVE_NW_BIT)
			}

			if tileMap.GetBit(row+1, col, EXTERIOR_WALL_N_BIT) && tileMap.GetBit(row, col-1, EXTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) {
				tileMap.SetBit(row+1, col-1, EXTERIOR_CORNER_CONCAVE_NE_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row-1, col+1, EXTERIOR_WALL_W_BIT) && tileMap.GetBit(row-1, col+1, EXTERIOR_WALL_S_BIT) {
				tileMap.SetBit(row-1, col+1, EXTERIOR_CORNER_CONVEX_SW_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row-1, col-1, EXTERIOR_WALL_E_BIT) && tileMap.GetBit(row-1, col-1, EXTERIOR_WALL_S_BIT) {
				tileMap.SetBit(row-1, col-1, EXTERIOR_CORNER_CONVEX_SE_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row+1, col+1, EXTERIOR_WALL_W_BIT) && tileMap.GetBit(row+1, col+1, EXTERIOR_WALL_N_BIT) {
				tileMap.SetBit(row+1, col+1, EXTERIOR_CORNER_CONVEX_NW_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row+1, col-1, EXTERIOR_WALL_E_BIT) && tileMap.GetBit(row+1, col-1, EXTERIOR_WALL_N_BIT) {
				tileMap.SetBit(row+1, col-1, EXTERIOR_CORNER_CONVEX_NE_BIT)
			}
		}
	}

	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col+1, INTERIOR_WALL_N_BIT) {
				setTerminus(tileMap, row, col)
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col-1, INTERIOR_WALL_N_BIT) {
				setTerminus(tileMap, row, col-1)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row, col+1, INTERIOR_WALL_S_BIT) {
				setTerminus(tileMap, row+1, col)
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row, col-1, INTERIOR_WALL_S_BIT) {
				setTerminus(tileMap, row+1, col-1)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row+1, col, INTERIOR_WALL_E_BIT) {
				setTerminus(tileMap, row+1, col)
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row-1, col, INTERIOR_WALL_E_BIT) {
				setTerminus(tileMap, row, col)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) && !tileMap.GetBit(row+1, col, INTERIOR_WALL_W_BIT) {
				setTerminus(tileMap, row+1, col-1)
			}
			if tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) && !tileMap.GetBit(row-1, col, INTERIOR_WALL_W_BIT) {
				setTerminus(tileMap, row, col-1)
			}
		}
	}
}

func setTerminus(tileMap *TileMap, row, col int) {
	tileMap.SetBit(row, col, TERMINUS_NE_BIT)
	tileMap.SetBit(row, col+1, TERMINUS_NW_BIT)
	tileMap.SetBit(row-1, col, TERMINUS_SE_BIT)
	tileMap.SetBit(row-1, col+1, TERMINUS_SW_BIT)
}

func (r *BaseRenderer) Render(x1, y1, x2, y2 float32, rooms []Room, navigationGraph *NavigationGraph, debugging bool) {
	tileMap := setAestheticBits(r.tileMap) // TODO - cache

	r.spriteBatch.Begin()

	startCol := math.Max(0, math.PrevMultiple(int(x1), 64)/64)
	startRow := math.Max(0, math.PrevMultiple(int(y1), 64)/64)

	endCol := math.Min(tileMap.Width(), math.NextMultiple(int(x2), 64)/64)
	endRow := math.Min(tileMap.Height(), math.NextMultiple(int(y2), 64)/64)

	// Floors
	for col := startCol; col < endCol; col++ {
		for row := startRow; row < endRow; row++ {
			for _, bitIndex := range TileBitIndexes {
				if bitIndex != FLOOR_BIT {
					continue
				}
				baseTexture, ok := r.textures[bitIndex]
				if !ok {
					continue
				}

				if tileMap.GetBit(row, col, bitIndex) {
					r.spriteBatch.Draw(
						baseTexture.texture,
						float32(col)*64, float32(row)*64, 64, 64,
						rendering.WithRotation(baseTexture.rotation),
						rendering.WithOrigin(32, 32),
					)
				}
			}
		}
	}

	// Fixtures
	for col := startCol; col < endCol; col++ {
		for row := startRow; row < endRow; row++ {
			if fixture, ok := tileMap.GetFixture(row, col); ok {
				r.spriteBatch.Draw(
					r.fixturesAtlas.Region(
						float32(fixture.AtlasRow)*64,
						float32(fixture.AtlasCol)*64,
						float32(fixture.TileWidth)*64,
						float32(fixture.TileHeight)*64,
					),
					float32(col)*64, float32(row)*64, float32(fixture.TileWidth)*64, float32(fixture.TileHeight)*64,
				)
			}
		}
	}

	// Non-floors
	for col := startCol; col < endCol; col++ {
		for row := startRow; row < endRow; row++ {
			for _, bitIndex := range TileBitIndexes {
				if bitIndex == FLOOR_BIT {
					continue
				}
				baseTexture, ok := r.textures[bitIndex]
				if !ok {
					continue
				}

				if tileMap.GetBit(row, col, bitIndex) {
					r.spriteBatch.Draw(
						baseTexture.texture,
						float32(col)*64, float32(row)*64, 64, 64,
						rendering.WithRotation(baseTexture.rotation),
						rendering.WithOrigin(32, 32),
					)
				}
			}
		}
	}

	if navigationGraph != nil && debugging {
		size := float32(10)
		lineSize := float32(2)

		for _, room := range rooms {
			for _, bound := range room.Bounds {
				c := room.Color
				c.A = 1
				for i, vertex := range bound.Vertices {
					r.spriteBatch.Draw(
						r.emptyTexture,
						vertex.X-size/2, vertex.Y-size/2, size, size,
						rendering.WithOrigin(size/2, size/2),
						rendering.WithColor(c),
					)

					to := bound.Vertices[(i+1)%len(bound.Vertices)]
					edge := to.Sub(vertex)
					angle := math.Atan232(edge.Y, edge.X)

					r.spriteBatch.Draw(
						r.emptyTexture,
						vertex.X, vertex.Y, edge.Len(), lineSize,
						rendering.WithRotation(angle),
						rendering.WithOrigin(0, 1),
						rendering.WithColor(c),
					)
				}
			}
		}

		if true {
			for _, node := range navigationGraph.Nodes {
				r.spriteBatch.Draw(
					r.emptyTexture,
					node.Center.X-size/2, node.Center.Y-size/2, size, size,
					rendering.WithOrigin(size/2, size/2),
					rendering.WithColor(rendering.Color{1, 0, 0, 1}),
				)
			}

			for _, edge := range navigationGraph.Edges {
				from := math.Vector{navigationGraph.Nodes[edge.From].Center.X - size/2, navigationGraph.Nodes[edge.From].Center.Y - size/2}
				to := math.Vector{navigationGraph.Nodes[edge.To].Center.X - size/2, navigationGraph.Nodes[edge.To].Center.Y - size/2}

				edge := to.Sub(from)
				angle := math.Atan232(edge.Y, edge.X)

				r.spriteBatch.Draw(
					r.emptyTexture,
					from.X+size/2, from.Y+size/2, edge.Len(), 1,
					rendering.WithRotation(angle),
					rendering.WithOrigin(0, 1),
					rendering.WithColor(rendering.Color{1, 0, 0, 1}),
				)
			}
		}
	}

	r.spriteBatch.End()
}

type baseTexture struct {
	texture  rendering.Texture
	rotation float32
}

func newBaseTexture(texture rendering.Texture, x, y int, rotation float32) baseTexture {
	return baseTexture{
		texture:  texture.Region(float32(x)*64, float32(y)*64, 64, 64),
		rotation: (rotation + 180) * stdmath.Pi / 180,
	}
}

func newBaseTextureMap(texture rendering.Texture, renderDoors bool) map[TileBitIndex]baseTexture {
	m := map[TileBitIndex]baseTexture{
		FLOOR_BIT:                      newBaseTexture(texture, 0, 1, 0),
		INTERIOR_WALL_N_BIT:            newBaseTexture(texture, 2, 1, 0),
		INTERIOR_WALL_S_BIT:            newBaseTexture(texture, 2, 1, 180),
		INTERIOR_WALL_E_BIT:            newBaseTexture(texture, 2, 1, 90),
		INTERIOR_WALL_W_BIT:            newBaseTexture(texture, 2, 1, 270),
		EXTERIOR_WALL_N_BIT:            newBaseTexture(texture, 1, 1, 0),
		EXTERIOR_WALL_S_BIT:            newBaseTexture(texture, 1, 1, 180),
		EXTERIOR_WALL_E_BIT:            newBaseTexture(texture, 1, 1, 90),
		EXTERIOR_WALL_W_BIT:            newBaseTexture(texture, 1, 1, 270),
		EXTERIOR_CORNER_CONVEX_NE_BIT:  newBaseTexture(texture, 1, 0, 90),
		EXTERIOR_CORNER_CONVEX_NW_BIT:  newBaseTexture(texture, 1, 0, 0),
		EXTERIOR_CORNER_CONVEX_SE_BIT:  newBaseTexture(texture, 1, 0, 180),
		EXTERIOR_CORNER_CONVEX_SW_BIT:  newBaseTexture(texture, 1, 0, 270),
		EXTERIOR_CORNER_CONCAVE_NE_BIT: newBaseTexture(texture, 2, 0, 90),
		EXTERIOR_CORNER_CONCAVE_NW_BIT: newBaseTexture(texture, 2, 0, 0),
		EXTERIOR_CORNER_CONCAVE_SE_BIT: newBaseTexture(texture, 2, 0, 180),
		EXTERIOR_CORNER_CONCAVE_SW_BIT: newBaseTexture(texture, 2, 0, 270),
		TERMINUS_SW_BIT:                newBaseTexture(texture, 0, 0, 270),
		TERMINUS_SE_BIT:                newBaseTexture(texture, 0, 0, 180),
		TERMINUS_NW_BIT:                newBaseTexture(texture, 0, 0, 0),
		TERMINUS_NE_BIT:                newBaseTexture(texture, 0, 0, 90),
	}

	if renderDoors {
		m[DOOR_N_BIT] = newBaseTexture(texture, 3, 1, 0)
		m[DOOR_S_BIT] = newBaseTexture(texture, 3, 1, 180)
		m[DOOR_E_BIT] = newBaseTexture(texture, 3, 1, 90)
		m[DOOR_W_BIT] = newBaseTexture(texture, 3, 1, 270)
	}

	return m
}
