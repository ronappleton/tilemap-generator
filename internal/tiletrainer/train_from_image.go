package tiletrainer

import (
	"fmt"
	"image"

	"tilemap-generator/internal/analyser"
	"tilemap-generator/internal/maputils"
)

func TrainFromImage(img image.Image, tileSize int, outputDir string, diagnostic bool) error {
	rawTiles := maputils.SliceImageIntoTiles(img, tileSize)
	groups, unique := analyser.FuzzyMatchTiles(rawTiles, 5)

	fmt.Printf("Deduplicated tiles: %d unique of %d total\n", unique, len(rawTiles))

	if diagnostic {
		diagPath := fmt.Sprintf("%s/diagnostic.png", outputDir)
		_ = analyser.SaveDiagnosticGrid(rawTiles, groups, tileSize, diagPath)
	}

	// Wrap into []Tile
	var tiles []maputils.Tile
	for i, tileImg := range rawTiles {
		hash, err := maputils.HashTile(tileImg)
		if err != nil {
			// skip tiles that fail to hash
			continue
		}
		tiles = append(tiles, maputils.Tile{
			ID:    i,
			Image: tileImg,
			Hash:  hash,
		})
	}

	return maputils.SaveTileset(tiles, outputDir, tileSize)
}
