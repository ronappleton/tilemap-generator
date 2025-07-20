package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"tilemap-generator/internal/analyser"
	"tilemap-generator/internal/iohelpers"
	"tilemap-generator/internal/tiletrainer"
)

var (
	inputName string
	tileSize  int
)

var trainTilesCmd = &cobra.Command{
	Use:   "train-tiles",
	Short: "Train tileset from a visual map",
	Run: func(cmd *cobra.Command, args []string) {
		resolvedPath, err := iohelpers.ResolveMapPath(inputName)
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}

		fmt.Println("🔍 Inspecting image...")
		analysis, err := analyser.InspectMap(resolvedPath)
		if err != nil {
			fmt.Println("❌ Failed to inspect image:", err)
			return
		}

		fmt.Printf("\n🗺️  Map Analysis:\n")
		fmt.Printf("- Resolution: %dx%d\n", analysis.Width, analysis.Height)
		fmt.Printf("- Colour Model: %s\n", analysis.ColorModel)
		fmt.Printf("- Unique Colours: %d\n", analysis.UniqueColors)
		fmt.Printf("- Uses Alpha Channel: %v\n", analysis.UsesAlpha)
		fmt.Printf("- Avg Brightness: %.1f\n", analysis.AvgBrightness)
		fmt.Printf("- Brightness Spread: %s\n", analysis.BrightnessSpread)

		fmt.Println("📊 Analysing image for optimal tile sizes...")
		candidateSizes := []int{16, 32, 64, 128, 256}
		results, err := analyser.AnalyseTileSizesFuzzy(resolvedPath, candidateSizes)
		if err != nil {
			fmt.Println("❌ Analysis failed:", err)
			return
		}

		fmt.Println("\nTile Size | Total Tiles | Unique Tiles | Reuse Ratio")
		fmt.Println("----------|-------------|---------------|-------------")
		for _, r := range results {
			fmt.Printf("%9d | %11d | %13d | %10.1f%%\n", r.TileSize, r.TotalTiles, r.UniqueTiles, r.ReuseRatio*100)
		}

		suggestedSize, ok := analyser.PickSuggestedTileSize(results, 0.3) // 30%+ reuse
		if !ok {
			fmt.Println("⚠️  No tile size offers sufficient reuse. Defaulting to 64.")
			suggestedSize = 64
		}

		fmt.Printf("\nSuggested tile size: %dpx — Proceed? (Y/n): ", suggestedSize)
		var answer string
		fmt.Scanln(&answer)
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer == "n" {
			fmt.Print("Enter custom tile size: ")
			var custom string
			fmt.Scanln(&custom)
			if i, err := strconv.Atoi(custom); err == nil {
				tileSize = i
			} else {
				fmt.Println("❌ Invalid number, aborting.")
				return
			}
		} else {
			tileSize = suggestedSize
		}

		baseName := strings.TrimSuffix(filepath.Base(resolvedPath), filepath.Ext(resolvedPath))
		outputDir := filepath.Join("tileset", baseName)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			fmt.Println("❌ Failed to create output directory:", err)
			return
		}

		fmt.Printf("\n🧠 Training tileset with %dpx tiles into '%s'...\n", tileSize, outputDir)
		// Load and clean the image
		img, err := imaging.Load(resolvedPath)
		if err != nil {
			fmt.Println("❌ Failed to load image:", err)
			return
		}
		cleaned := analyser.PreprocessForTraining(img)
		tiletrainer.Train(resolvedPath, tileSize, outputDir)
	},
}

func init() {
	trainTilesCmd.Flags().StringVarP(&inputName, "input", "i", "", "Name of map to train on (without extension)")
	trainTilesCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(trainTilesCmd)
}
