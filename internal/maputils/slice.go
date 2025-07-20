package maputils

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Tile struct {
	ID    int
	Image image.Image
	Hash  string
	X     int
	Y     int
}

func SliceAndHashTiles(path string, tileSize int) ([]Tile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	bounds := img.Bounds()
	cols := bounds.Dx() / tileSize
	rows := bounds.Dy() / tileSize

	var tiles []Tile
	tileID := 0

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			tileRect := image.Rect(x*tileSize, y*tileSize, (x+1)*tileSize, (y+1)*tileSize)
			subImage := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(tileRect)

			// TODO: Add deduplication logic here using hash
			tiles = append(tiles, Tile{
				ID:    tileID,
				Image: subImage,
				Hash:  "", // placeholder
				X:     x,
				Y:     y,
			})
			tileID++
		}
	}

	return tiles, nil
}
