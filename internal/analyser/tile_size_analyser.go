package analyser

import (
	"fmt"
	"image"
	"os"
	"sort"

	"tilemap-generator/internal/maputils"
)

type TileSizeResult struct {
	TileSize    int
	TotalTiles  int
	UniqueTiles int
	ReuseRatio  float64
}

func AnalyseTileSizes(imagePath string, candidateSizes []int) ([]TileSizeResult, error) {
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image: %v", err)
	}

	// Clean the image before analysis
	img = PreprocessForTraining(img)

	var results []TileSizeResult
	for _, size := range candidateSizes {
		tiles := maputils.SliceImageIntoTiles(img, size)
		hashes := maputils.HashTiles(tiles)
		unique := maputils.DeduplicateTiles(hashes)

		total := len(tiles)
		uniqueCount := len(unique)
		ratio := float64(total-uniqueCount) / float64(total)

		results = append(results, TileSizeResult{
			TileSize:    size,
			TotalTiles:  total,
			UniqueTiles: uniqueCount,
			ReuseRatio:  ratio,
		})
	}

	// Sort by TileSize ascending
	sort.Slice(results, func(i, j int) bool {
		return results[i].TileSize < results[j].TileSize
	})

	return results, nil
}

func PickSuggestedTileSize(results []TileSizeResult, minReuse float64) (int, bool) {
	for _, result := range results {
		if result.ReuseRatio >= minReuse {
			return result.TileSize, true
		}
	}
	return 0, false
}
