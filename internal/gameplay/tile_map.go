package gameplay

import (
	"bufio"
	"encoding/binary"
	"io"
)

type TileBitIndex int64

const (
	// TODO - rename
	FLOOR_BIT TileBitIndex = iota
	HWALL_BIT
	VWALL_BIT
	HDOOR_BIT
	VDOOR_BIT
	INTERIOR_WALL_N_BIT
	INTERIOR_WALL_S_BIT
	INTERIOR_WALL_E_BIT
	INTERIOR_WALL_W_BIT
	EXTERIOR_WALL_N_BIT
	EXTERIOR_WALL_S_BIT
	EXTERIOR_WALL_E_BIT
	EXTERIOR_WALL_W_BIT
	EXTERIOR_CORNER_CONVEX_NW_BIT
	EXTERIOR_CORNER_CONVEX_NE_BIT
	EXTERIOR_CORNER_CONVEX_SW_BIT
	EXTERIOR_CORNER_CONVEX_SE_BIT
	EXTERIOR_CORNER_CONCAVE_NW_BIT
	EXTERIOR_CORNER_CONCAVE_NE_BIT
	EXTERIOR_CORNER_CONCAVE_SE_BIT
	EXTERIOR_CORNER_CONCAVE_SW_BIT
	TERMINUS_SW_BIT
	TERMINUS_NW_BIT
	TERMINUS_SE_BIT
	TERMINUS_NE_BIT
	DOOR_N_BIT
	DOOR_E_BIT
	DOOR_S_BIT
	DOOR_W_BIT
)

var tileBitIndexes = []TileBitIndex{
	FLOOR_BIT,
	HWALL_BIT,
	VWALL_BIT,
	HDOOR_BIT,
	VDOOR_BIT,
	INTERIOR_WALL_N_BIT,
	INTERIOR_WALL_S_BIT,
	INTERIOR_WALL_E_BIT,
	INTERIOR_WALL_W_BIT,
	EXTERIOR_WALL_N_BIT,
	EXTERIOR_WALL_S_BIT,
	EXTERIOR_WALL_E_BIT,
	EXTERIOR_WALL_W_BIT,
	EXTERIOR_CORNER_CONVEX_NW_BIT,
	EXTERIOR_CORNER_CONVEX_NE_BIT,
	EXTERIOR_CORNER_CONVEX_SW_BIT,
	EXTERIOR_CORNER_CONVEX_SE_BIT,
	EXTERIOR_CORNER_CONCAVE_NW_BIT,
	EXTERIOR_CORNER_CONCAVE_NE_BIT,
	EXTERIOR_CORNER_CONCAVE_SE_BIT,
	EXTERIOR_CORNER_CONCAVE_SW_BIT,
	TERMINUS_SW_BIT,
	TERMINUS_NW_BIT,
	TERMINUS_SE_BIT,
	TERMINUS_NE_BIT,
	DOOR_N_BIT,
	DOOR_E_BIT,
	DOOR_S_BIT,
	DOOR_W_BIT,
}

type TileMap struct {
	width    int
	height   int
	gridSize int
	data     []int64
}

func ReadTileMap(r io.Reader) (*TileMap, error) {
	byteReader := bufio.NewReader(r)

	width, err := binary.ReadVarint(byteReader)
	if err != nil {
		return nil, err
	}

	height, err := binary.ReadVarint(byteReader)
	if err != nil {
		return nil, err
	}

	gridSize, err := binary.ReadVarint(byteReader)
	if err != nil {
		return nil, err
	}

	m := NewTileMap(int(width), int(height), int(gridSize))
	for i := 0; i < len(m.data); i++ {
		val, err := binary.ReadVarint(byteReader)
		if err != nil {
			return nil, err
		}

		m.data[i] = val
	}

	return m, nil
}

func WriteTileMap(m *TileMap, w io.Writer) error {
	var buf []byte
	buf = binary.AppendVarint(buf, int64(m.width))
	buf = binary.AppendVarint(buf, int64(m.height))
	buf = binary.AppendVarint(buf, int64(m.gridSize))

	for _, val := range m.data {
		buf = binary.AppendVarint(buf, int64(val))
	}

	_, err := w.Write(buf)
	return err
}

func NewTileMap(width, height, gridSize int) *TileMap {
	return &TileMap{
		width:    width,
		height:   height,
		gridSize: gridSize,
		data:     make([]int64, width*height),
	}
}

func (m *TileMap) GetBit(row, col int, bitIndex TileBitIndex) bool {
	return m.GetAllBits(row, col, bitIndex)
}

func (m *TileMap) GetAllBits(row, col int, bitIndexes ...TileBitIndex) bool {
	bits := m.GetBits(row, col)

	for _, bit := range bitIndexes {
		if bits&(1<<bit) == 0 {
			return false
		}
	}

	return true
}

func (m *TileMap) GetBits(row, col int) int64 {
	return m.data[m.bitsetIndex(row, col)]
}

func (m *TileMap) SetBit(row, col int, bitIndex TileBitIndex) {
	m.SetBits(row, col, m.GetBits(row, col)|(1<<bitIndex))
}

func (m *TileMap) SetBits(row, col int, val int64) {
	m.data[m.bitsetIndex(row, col)] = val
}

func (m *TileMap) ClearBit(row, col int, bitIndex TileBitIndex) {
	m.SetBits(row, col, m.GetBits(row, col) & ^(1<<bitIndex))
}

func (m *TileMap) ClearBits(row, col int) {
	m.SetBits(row, col, 0)
}

func (m *TileMap) ClearAll() {
	for i := range m.data {
		m.data[i] = 0
	}
}

func (m *TileMap) bitsetIndex(row, col int) int {
	return col*m.width + row
}
