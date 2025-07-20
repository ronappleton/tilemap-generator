package tiletrainer

import (
	"fmt"

	"tilemap-generator/internal/imagehelpers"
	"tilemap-generator/internal/tileutils"
)

func Train(imagePath string, tileSize int, outputDir string) error {
	orig, cleaned, err := imagehelpers.LoadAndPreprocess(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	tiles, mapping, err := tileutils.ExtractUniqueTilesWithIndex(orig, cleaned, tileSize)
	if err != nil {
		return fmt.Errorf("failed to extract tiles: %w", err)
	}

	return tileutils.SaveTilesetWithIndex(tiles, mapping, outputDir, tileSize)
}
