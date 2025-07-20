package imagehelpers

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

// LoadAndPreprocess opens an image from disk and returns both the original image
// and a preprocessed version for tile analysis.
func LoadAndPreprocess(path string) (image.Image, image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open image: %w", err)
	}
	cleaned := PreprocessForTraining(img)
	return img, cleaned, nil
}
