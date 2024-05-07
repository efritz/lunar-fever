package maps

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

func NewBaseRenderSystem(engineCtx *engine.Context, tileMap *TileMap) *baseRenderSystem {
	return &baseRenderSystem{
		Context: engineCtx,
		tileMap: tileMap,
	}
}

func (s *baseRenderSystem) Init() {
	s.baseRenderer = NewBaseRenderer(s.SpriteBatch, s.TextureLoader, s.tileMap, false)
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

func NewBaseRenderer(spriteBatch *rendering.SpriteBatch, textureLoader *rendering.TextureLoader, tileMap *TileMap, renderDoors bool) *BaseRenderer {
	return &BaseRenderer{
		spriteBatch: spriteBatch,
		tileMap:     tileMap,
		textures:    newBaseTextureMap(textureLoader.Load("base"), renderDoors),
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
	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
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

	for col := 0; col < tileMap.Width(); col++ {
		for row := 0; row < tileMap.Height(); row++ {
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
	tileMap := setAestheticBits(r.tileMap) // TODO - cache

	r.spriteBatch.Begin()

	startCol := math.Max(0, math.PrevMultiple(int(x1), 64)/64)
	startRow := math.Max(0, math.PrevMultiple(int(y1), 64)/64)

	endCol := math.Min(tileMap.Width(), math.NextMultiple(int(x2), 64)/64)
	endRow := math.Min(tileMap.Height(), math.NextMultiple(int(y2), 64)/64)

	for col := startCol; col < endCol; col++ {
		for row := startRow; row < endRow; row++ {
			for _, bitIndex := range TileBitIndexes {
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
