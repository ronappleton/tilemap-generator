package tileutils

import (
	"image"

	"tilemap-generator/internal/maputils"
)

// ExtractUniqueTilesWithIndex slices both original and cleaned images into tiles,
// deduplicates tiles using the cleaned version, and returns unique tiles from the
// original along with a mapping of tile indices to unique tile IDs.
func ExtractUniqueTilesWithIndex(original, cleaned image.Image, tileSize int) ([]maputils.Tile, [][]int, error) {
	cleanTiles := maputils.SliceImageIntoTiles(cleaned, tileSize)
	origTiles := maputils.SliceImageIntoTiles(original, tileSize)

	if len(cleanTiles) != len(origTiles) {
		return nil, nil, nil
	}

	bounds := cleaned.Bounds()
	cols := bounds.Dx() / tileSize
	rows := bounds.Dy() / tileSize

	mapping := make([][]int, rows)
	for i := range mapping {
		mapping[i] = make([]int, cols)
	}

	seen := make(map[string]int)
	var tiles []maputils.Tile
	nextID := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			idx := y*cols + x
			hash, err := maputils.HashTile(cleanTiles[idx])
			if err != nil {
				return nil, nil, err
			}
			id, ok := seen[hash]
			if !ok {
				id = nextID
				seen[hash] = id
				tiles = append(tiles, maputils.Tile{
					ID:    id,
					Image: origTiles[idx],
					Hash:  hash,
					X:     x,
					Y:     y,
				})
				nextID++
			}
			mapping[y][x] = id
		}
	}

	return tiles, mapping, nil
}

// ExtractTiles returns a slice of tiles from an image.
func ExtractTiles(img image.Image, tileSize int) []image.Image {
	return maputils.SliceImageIntoTiles(img, tileSize)
}
