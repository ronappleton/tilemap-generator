package analyser

import (
	"image"
	_ "image/color"
	"math"

	"github.com/disintegration/gift"
)

// PreprocessForTraining takes an image and returns a cleaned version suitable for tile analysis.
// It removes anti-aliasing, reduces colours, and enforces hard edges.
func PreprocessForTraining(img image.Image) image.Image {
	g := gift.New(
		// Step 1: Slight blur to smooth noise (anti-aliasing removal)
		gift.GaussianBlur(0.4),
	)

	// Create a destination image and apply the filter pipeline
	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	// Step 2: Manual posterize
	posterizeImage(dst, 4) // 4 levels per channel

	// Step 3: Enforce hard edges via per-pixel thresholding
	applyThreshold(dst, 15)

	return dst
}

// applyThreshold walks each pixel and clamps it to fully dark or light if close to edge.
func applyThreshold(img *image.RGBA, threshold uint8) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			o := img.PixOffset(x, y)
			r, g, b := img.Pix[o], img.Pix[o+1], img.Pix[o+2]

			// Simplistic thresholding logic
			avg := uint8((uint16(r) + uint16(g) + uint16(b)) / 3)
			if math.Abs(float64(avg)-128) < float64(threshold) {
				img.Pix[o] = 128
				img.Pix[o+1] = 128
				img.Pix[o+2] = 128
			}
		}
	}
}

// posterizeImage reduces the number of colors by snapping RGB values to fixed levels.
func posterizeImage(img *image.RGBA, levels int) {
	step := 256 / levels
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			o := img.PixOffset(x, y)
			img.Pix[o+0] = uint8((int(img.Pix[o+0]) / step) * step)
			img.Pix[o+1] = uint8((int(img.Pix[o+1]) / step) * step)
			img.Pix[o+2] = uint8((int(img.Pix[o+2]) / step) * step)
		}
	}
}
