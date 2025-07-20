package tiletrainer

import (
	"fmt"
	"tilemap-generator/internal/maputils"
)

func Train(imagePath string, tileSize int, outputDir string) error {
	tiles, err := maputils.SliceAndHashTiles(imagePath, tileSize)
	if err != nil {
		return fmt.Errorf("failed to process tiles: %w", err)
	}

	err = maputils.SaveTileset(tiles, outputDir, tileSize) // ← Fix: pass tileSize here
	if err != nil {
		return fmt.Errorf("failed to save tileset: %w", err)
	}

	return nil
}
