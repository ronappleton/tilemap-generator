package tiletrainer

import (
	"fmt"
	"image"

	"tilemap-generator/internal/analyser"
	"tilemap-generator/internal/tileutils"
)

// TrainFromImages deduplicates tiles using the cleaned image but saves tiles
// cut from the original image into outputDir. A mapping of tile positions to
// tile IDs is written to tileset.json.
func TrainFromImages(original, cleaned image.Image, tileSize int, outputDir string, diagnostic bool) error {
	rawTiles := tileutils.ExtractTiles(cleaned, tileSize)
	groups, unique := analyser.FuzzyMatchTiles(rawTiles, 5)

	fmt.Printf("Deduplicated tiles: %d unique of %d total\n", unique, len(rawTiles))

	if diagnostic {
		diagPath := fmt.Sprintf("%s/diagnostic.png", outputDir)
		_ = analyser.SaveDiagnosticGrid(rawTiles, groups, tileSize, diagPath)
	}

	tiles, mapping, err := tileutils.ExtractUniqueTilesWithIndex(original, cleaned, tileSize)
	if err != nil {
		return err
	}

	return tileutils.SaveTilesetWithIndex(tiles, mapping, outputDir, tileSize)
}
