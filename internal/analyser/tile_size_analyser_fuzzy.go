// internal/analyser/tile_size_analyser_fuzzy.go
package analyser

import (
	"fmt"
	"image"
	"os"

	"tilemap-generator/internal/maputils"
)

func AnalyseTileSizesFuzzy(imgPath string, sizes []int) ([]TileSizeResult, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Clean the image before analysis
	srcImg = PreprocessForTraining(srcImg)

	var results []TileSizeResult
	for _, size := range sizes {
		tiles := maputils.SliceImageIntoTiles(srcImg, size)
		_, unique := FuzzyMatchTiles(tiles, 5)

		reuseRatio := 1.0 - float64(unique)/float64(len(tiles))
		results = append(results, TileSizeResult{
			TileSize:    size,
			TotalTiles:  len(tiles),
			UniqueTiles: unique,
			ReuseRatio:  reuseRatio,
		})
	}

	return results, nil
}
