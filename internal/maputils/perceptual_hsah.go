package maputils

import (
	"fmt"
	"image"

	"github.com/corona10/goimagehash"
)

func PerceptualHash(img image.Image) (string, error) {
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", fmt.Errorf("failed to generate pHash: %w", err)
	}
	return fmt.Sprintf("%d", hash.GetHash()), nil
}
