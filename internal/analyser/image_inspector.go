package analyser

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"sort"
)

type MapAnalysisResult struct {
	Width            int
	Height           int
	ColorModel       string
	UniqueColors     int
	UsesAlpha        bool
	AvgBrightness    float64
	BrightnessSpread string
}

func InspectMap(imgPath string) (*MapAnalysisResult, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	decoded, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Clean the image before collecting statistics
	decoded = PreprocessForTraining(decoded)

	bounds := decoded.Bounds()
	uniqueColors := make(map[color.Color]struct{})
	brightnessValues := []float64{}
	usesAlpha := false
	var totalBrightness float64

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := decoded.At(x, y)
			r, g, b, a := c.RGBA()

			if a < 0xffff {
				usesAlpha = true
			}

			norm := func(v uint32) float64 { return float64(v) / 65535.0 }
			brightness := (norm(r) + norm(g) + norm(b)) / 3
			brightnessValues = append(brightnessValues, brightness)
			totalBrightness += brightness

			uniqueColors[c] = struct{}{}
		}
	}

	sort.Float64s(brightnessValues)
	min := brightnessValues[0]
	max := brightnessValues[len(brightnessValues)-1]
	avg := totalBrightness / float64(len(brightnessValues))

	analysis := &MapAnalysisResult{
		Width:            bounds.Dx(),
		Height:           bounds.Dy(),
		ColorModel:       fmt.Sprintf("%T", decoded.ColorModel()),
		UniqueColors:     len(uniqueColors),
		UsesAlpha:        usesAlpha,
		AvgBrightness:    avg * 255,
		BrightnessSpread: fmt.Sprintf("%.0f–%.0f", min*255, max*255),
	}

	return analysis, nil
}
