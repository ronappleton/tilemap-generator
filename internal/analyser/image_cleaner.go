package analyser

import (
	"image"

	"tilemap-generator/internal/imagehelpers"
)

// PreprocessForTraining takes an image and returns a cleaned version suitable for tile analysis.
// It delegates to imagehelpers.PreprocessForTraining.
func PreprocessForTraining(img image.Image) image.Image {
	return imagehelpers.PreprocessForTraining(img)
}
