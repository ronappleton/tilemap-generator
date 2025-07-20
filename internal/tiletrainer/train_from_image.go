package tiletrainer

import (
	"image"
	"tilemap-generator/internal/maputils"
)

func TrainFromImage(img image.Image, tileSize int, outputDir string) {
	rawTiles := maputils.SliceImageIntoTiles(img, tileSize)

	// Wrap into []Tile
	var tiles []maputils.Tile
	for i, tileImg := range rawTiles {
		tiles = append(tiles, maputils.Tile{
			ID:   i,
			Img:  tileImg,
			Hash: maputils.HashImage(tileImg),
		})
	}

	maputils.SaveTileset(tiles, outputDir, tileSize)
}
