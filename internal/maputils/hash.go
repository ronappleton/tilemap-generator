package maputils

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"image"
	"image/png"
)

// HashTile creates a SHA-1 hash of a tile's PNG-encoded byte representation.
func HashTile(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	sum := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(sum[:]), nil
}
