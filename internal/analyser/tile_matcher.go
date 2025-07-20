package analyser

import (
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/bits"
	"os"

	"tilemap-generator/internal/maputils"
)

// FuzzyMatchTiles groups tiles by similarity using aHash and a Hamming distance threshold.
// It returns a slice assigning each tile to a group and the number of unique groups.
func FuzzyMatchTiles(tiles []image.Image, threshold int) ([]int, int) {
	var uniqueHashes []string
	assignments := make([]int, len(tiles))

	for i, t := range tiles {
		h := maputils.FuzzyHash(t)
		matched := false
		for id, u := range uniqueHashes {
			if hammingDistance(h, u) <= threshold {
				assignments[i] = id
				matched = true
				break
			}
		}
		if !matched {
			assignments[i] = len(uniqueHashes)
			uniqueHashes = append(uniqueHashes, h)
		}
	}

	return assignments, len(uniqueHashes)
}

func hammingDistance(a, b string) int {
	ab, _ := hex.DecodeString(a)
	bb, _ := hex.DecodeString(b)
	d := 0
	for i := 0; i < len(ab) && i < len(bb); i++ {
		d += bits.OnesCount8(ab[i] ^ bb[i])
	}
	return d
}

// SaveDiagnosticGrid creates a PNG showing all tiles with coloured borders for their groups.
func SaveDiagnosticGrid(tiles []image.Image, groups []int, tileSize int, path string) error {
	if len(tiles) == 0 {
		return nil
	}
	cols := int(math.Ceil(math.Sqrt(float64(len(tiles)))))
	rows := int(math.Ceil(float64(len(tiles)) / float64(cols)))

	outImg := image.NewRGBA(image.Rect(0, 0, cols*tileSize, rows*tileSize))
	for idx, t := range tiles {
		x := (idx % cols) * tileSize
		y := (idx / cols) * tileSize
		r := image.Rect(x, y, x+tileSize, y+tileSize)
		draw.Draw(outImg, r, t, image.Point{}, draw.Src)
		drawBorder(outImg, r, colourForGroup(groups[idx]))
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, outImg)
}

func drawBorder(img *image.RGBA, r image.Rectangle, c color.Color) {
	for x := r.Min.X; x < r.Max.X; x++ {
		img.Set(x, r.Min.Y, c)
		img.Set(x, r.Max.Y-1, c)
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		img.Set(r.Min.X, y, c)
		img.Set(r.Max.X-1, y, c)
	}
}

func colourForGroup(id int) color.Color {
	return color.RGBA{uint8((id * 53) % 256), uint8((id * 97) % 256), uint8((id * 191) % 256), 255}
}
