package loader

import (
	"os"

	"github.com/efritz/lunar-fever/internal/gameplay/maps"
)

const tempPath = "map.dat"

func ReadTileMap() (*maps.TileMap, error) {
	f, err := os.Open(tempPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return maps.ReadTileMap(f)
}

func Write(tileMap *maps.TileMap) error {
	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := maps.WriteTileMap(tileMap, f); err != nil {
		return err
	}

	return nil
}
