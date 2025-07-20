package tiletrainer

import (
	"image"
	"tilemap-generator/internal/maputils"
)

func TrainFromImage(img image.Image, tileSize int, outputDir string) error {
	rawTiles := maputils.SliceImageIntoTiles(img, tileSize)

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
