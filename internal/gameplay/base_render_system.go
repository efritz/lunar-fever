package gameplay

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type baseRenderSystem struct {
	*engine.Context
	tileMap      *TileMap
	baseRenderer *BaseRenderer
}

func (s *baseRenderSystem) Init() {
	s.baseRenderer = NewBaseRenderer(s.SpriteBatch, s.TextureLoader, s.tileMap)
}

func (s *baseRenderSystem) Exit() {}

func (s *baseRenderSystem) Process(elapsedMs int64) {
	s.baseRenderer.Render(s.Camera.Bounds())
}

//
//

type BaseRenderer struct {
	spriteBatch *rendering.SpriteBatch
	tileMap     *TileMap
	textures    map[TileBitIndex]baseTexture
}

func NewBaseRenderer(spriteBatch *rendering.SpriteBatch, textureLoader *rendering.TextureLoader, tileMap *TileMap) *BaseRenderer {
	return &BaseRenderer{
		spriteBatch: spriteBatch,
		tileMap:     tileMap,
		textures:    newBaseTextureMap(textureLoader.Load("base")),
	}
}

func setAestheticBits(tileMap *TileMap) *TileMap {
	m2 := NewTileMap(tileMap.width, tileMap.height, tileMap.gridSize)
	for col := 0; col < tileMap.width; col++ {
		for row := 0; row < tileMap.height; row++ {
			m2.SetBits(row, col, tileMap.GetBits(row, col))
		}
	}

	setExternalWallTiles(m2)
	setCornerTiles(m2)
	return m2
}

func setExternalWallTiles(tileMap *TileMap) {
	for col := 0; col < tileMap.width; col++ {
		for row := 0; row < tileMap.height; row++ {
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row-1, col, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_WALL_N_BIT)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_S_BIT) && !tileMap.GetBit(row+1, col, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_WALL_S_BIT)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_WALL_E_BIT)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_W_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_WALL_W_BIT)
			}
		}
	}
}

func setCornerTiles(tileMap *TileMap) {
	for col := 0; col < tileMap.width; col++ {
		for row := 0; row < tileMap.height; row++ {
			if tileMap.GetBit(row, col, EXTERIOR_WALL_N_BIT) && tileMap.GetBit(row, col, EXTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONVEX_NE_BIT)
			}

			if tileMap.GetBit(row, col, EXTERIOR_WALL_N_BIT) && tileMap.GetBit(row, col, EXTERIOR_WALL_W_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONVEX_NW_BIT)
			}

			if tileMap.GetBit(row, col, EXTERIOR_WALL_S_BIT) && tileMap.GetBit(row, col, EXTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONVEX_SE_BIT)
			}

			if tileMap.GetBit(row, col, EXTERIOR_WALL_S_BIT) && tileMap.GetBit(row, col, EXTERIOR_WALL_W_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONVEX_SW_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row-1, col, EXTERIOR_WALL_E_BIT) && tileMap.GetBit(row, col+1, EXTERIOR_WALL_N_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONCAVE_NE_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row-1, col, EXTERIOR_WALL_W_BIT) && tileMap.GetBit(row, col-1, EXTERIOR_WALL_N_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONCAVE_NW_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row+1, col, EXTERIOR_WALL_E_BIT) && tileMap.GetBit(row, col+1, EXTERIOR_WALL_S_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONCAVE_SE_BIT)
			}

			if tileMap.GetBit(row, col, FLOOR_BIT) && tileMap.GetBit(row+1, col, EXTERIOR_WALL_W_BIT) && tileMap.GetBit(row, col-1, EXTERIOR_WALL_S_BIT) {
				tileMap.SetBit(row, col, EXTERIOR_CORNER_CONCAVE_SW_BIT)
			}
		}
	}

	for col := 0; col < tileMap.width; col++ {
		for row := 0; row < tileMap.height; row++ {
			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col+1, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col, EXTERIOR_WALL_E_BIT) {
				setTerminus(tileMap, row, col)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col-1, INTERIOR_WALL_N_BIT) && !tileMap.GetBit(row, col, EXTERIOR_WALL_W_BIT) {
				setTerminus(tileMap, row, col-1)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row-1, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col, EXTERIOR_WALL_N_BIT) {
				setTerminus(tileMap, row, col)
			}

			if tileMap.GetBit(row, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row+1, col, INTERIOR_WALL_E_BIT) && !tileMap.GetBit(row, col, EXTERIOR_WALL_S_BIT) {
				setTerminus(tileMap, row+1, col)
			}
		}
	}
}

func setTerminus(tileMap *TileMap, row, col int) {
	if tileMap.GetBit(row, col, EXTERIOR_CORNER_CONVEX_NE_BIT) || tileMap.GetBit(row, col, EXTERIOR_CORNER_CONCAVE_NE_BIT) {
		return
	}
	if tileMap.GetBit(row, col+1, EXTERIOR_CORNER_CONVEX_NW_BIT) || tileMap.GetBit(row, col+1, EXTERIOR_CORNER_CONCAVE_NW_BIT) {
		return
	}
	if tileMap.GetBit(row-1, col, EXTERIOR_CORNER_CONVEX_SE_BIT) || tileMap.GetBit(row-1, col, EXTERIOR_CORNER_CONCAVE_SE_BIT) {
		return
	}
	if tileMap.GetBit(row-1, col+1, EXTERIOR_CORNER_CONVEX_SW_BIT) || tileMap.GetBit(row-1, col+1, EXTERIOR_CORNER_CONCAVE_SW_BIT) {
		return
	}

	tileMap.SetBit(row, col, TERMINUS_NE_BIT)
	tileMap.SetBit(row, col+1, TERMINUS_NW_BIT)
	tileMap.SetBit(row-1, col, TERMINUS_SE_BIT)
	tileMap.SetBit(row-1, col+1, TERMINUS_SW_BIT)
}

func (r *BaseRenderer) Render(x1, y1, x2, y2 float32) {
	tileMap := setAestheticBits(r.tileMap)

	r.spriteBatch.Begin()

	startCol := math.Max(0, math.PrevMultiple(int(x1), 64)/64)
	startRow := math.Max(0, math.PrevMultiple(int(y1), 64)/64)

	endCol := math.Min(tileMap.width, math.NextMultiple(int(x2), 64)/64)
	endRow := math.Min(tileMap.height, math.NextMultiple(int(y2), 64)/64)

	for col := startCol; col < endCol; col++ {
		for row := startRow; row < endRow; row++ {
			for _, bitIndex := range tileBitIndexes {
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

func newBaseTextureMap(texture rendering.Texture) map[TileBitIndex]baseTexture {
	return map[TileBitIndex]baseTexture{
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
		DOOR_N_BIT:                     newBaseTexture(texture, 3, 1, 0),
		DOOR_S_BIT:                     newBaseTexture(texture, 3, 1, 180),
		DOOR_E_BIT:                     newBaseTexture(texture, 3, 1, 90),
		DOOR_W_BIT:                     newBaseTexture(texture, 3, 1, 270),
		TERMINUS_SW_BIT:                newBaseTexture(texture, 0, 0, 270),
		TERMINUS_SE_BIT:                newBaseTexture(texture, 0, 0, 180),
		TERMINUS_NW_BIT:                newBaseTexture(texture, 0, 0, 0),
		TERMINUS_NE_BIT:                newBaseTexture(texture, 0, 0, 90),
	}
}
