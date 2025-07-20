package analyser

import (
	"image"
	"math"

	"github.com/disintegration/gift"
)

// PreprocessForTraining takes an image and returns a cleaned version suitable for tile analysis.
// It removes anti-aliasing, reduces colours, and enforces hard edges.
func PreprocessForTraining(img image.Image) image.Image {
	// Grayscale conversion ensures consistent processing and reduces colour noise
	g := gift.New(
		gift.Grayscale(),
		// Slight blur to smooth noise (anti-aliasing removal)
		gift.GaussianBlur(0.4),
		// Sharpen edges after blur to enhance boundaries
		gift.UnsharpMask(1, 1, 0),
	)

	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	// First posterisation to flatten colours
	posterizeImage(dst, 4)
	// Morphological opening to remove isolated noise pixels
	morphologicalOpen(dst)
	// Reapply posterisation in case sharpening introduced small variations
	posterizeImage(dst, 4)

	// Enforce hard edges via per-pixel thresholding
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

// morphologicalOpen performs an erosion followed by dilation using a 3x3 kernel.
func morphologicalOpen(img *image.RGBA) {
	eroded := erode(img)
	dilated := dilate(eroded)
	copy(img.Pix, dilated.Pix)
}

func erode(src *image.RGBA) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			min := uint8(255)
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky
					if px < b.Min.X || py < b.Min.Y || px >= b.Max.X || py >= b.Max.Y {
						continue
					}
					o := src.PixOffset(px, py)
					if src.Pix[o] < min {
						min = src.Pix[o]
					}
				}
			}
			o := dst.PixOffset(x, y)
			dst.Pix[o] = min
			dst.Pix[o+1] = min
			dst.Pix[o+2] = min
			dst.Pix[o+3] = src.Pix[o+3]
		}
	}
	return dst
}

func dilate(src *image.RGBA) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			max := uint8(0)
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky
					if px < b.Min.X || py < b.Min.Y || px >= b.Max.X || py >= b.Max.Y {
						continue
					}
					o := src.PixOffset(px, py)
					if src.Pix[o] > max {
						max = src.Pix[o]
					}
				}
			}
			o := dst.PixOffset(x, y)
			dst.Pix[o] = max
			dst.Pix[o+1] = max
			dst.Pix[o+2] = max
			dst.Pix[o+3] = src.Pix[o+3]
		}
	}
	return dst
}
