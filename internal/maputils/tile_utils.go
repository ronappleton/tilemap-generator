package maputils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/draw"
	"image/png"
)

// SliceImageIntoTiles slices an image into square tiles of a given size.
func SliceImageIntoTiles(img image.Image, tileSize int) []image.Image {
	bounds := img.Bounds()
	var tiles []image.Image

	for y := bounds.Min.Y; y+tileSize <= bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x+tileSize <= bounds.Max.X; x += tileSize {
			tileRect := image.Rect(0, 0, tileSize, tileSize)
			tile := image.NewRGBA(tileRect)
			draw.Draw(tile, tileRect, img, image.Point{X: x, Y: y}, draw.Src)
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

// HashTiles returns an array of string hashes representing each tile.
func HashTiles(tiles []image.Image) []string {
	var hashes []string
	for _, tile := range tiles {
		buf := new(bytes.Buffer)
		h := md5.New()
		err := png.Encode(buf, tile)
		if err != nil {
			continue
		}
		h.Write(buf.Bytes())
		hashes = append(hashes, hex.EncodeToString(h.Sum(nil)))
	}
	return hashes
}

// DeduplicateTiles returns a map of unique hash -> count of occurrences.
func DeduplicateTiles(hashes []string) map[string]int {
	hashMap := make(map[string]int)
	for _, h := range hashes {
		hashMap[h]++
	}
	return hashMap
}
