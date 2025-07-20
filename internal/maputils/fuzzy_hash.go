package maputils

import (
	"bytes"
	_ "encoding/binary"
	"image"
	_ "image/color"
	_ "math"

	"golang.org/x/image/draw"
)

func FuzzyHash(img image.Image) string {
	// Downscale to 8x8 grayscale using NearestNeighbor
	thumb := image.NewGray(image.Rect(0, 0, 8, 8))
	draw.NearestNeighbor.Scale(thumb, thumb.Bounds(), img, img.Bounds(), draw.Over, nil)

	var hash bytes.Buffer
	avg := averageGray(thumb)

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			c := thumb.GrayAt(x, y).Y
			if c > avg {
				hash.WriteByte(1)
			} else {
				hash.WriteByte(0)
			}
		}
	}

	return binaryToHex(hash.Bytes())
}

func averageGray(img *image.Gray) uint8 {
	var total uint64
	for _, val := range img.Pix {
		total += uint64(val)
	}
	return uint8(total / uint64(len(img.Pix)))
}

func binaryToHex(data []byte) string {
	var result string
	for i := 0; i < len(data); i += 8 {
		var b byte
		for j := 0; j < 8 && i+j < len(data); j++ {
			b |= (data[i+j] & 1) << (7 - j)
		}
		result += string("0123456789abcdef"[b>>4])
		result += string("0123456789abcdef"[b&0x0F])
	}
	return result
}
