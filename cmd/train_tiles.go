package cmd

import (
	"fmt"
	"os"
	_ "path/filepath"

	"tilemap-generator/internal/iohelpers"
	"tilemap-generator/internal/tiletrainer"

	"github.com/spf13/cobra"
)

var (
	inputName string
	tileSize  int
	outputDir string
)

var trainTilesCmd = &cobra.Command{
	Use:   "train-tiles",
	Short: "Slice and learn tile rules from a PNG map",
	Run: func(cmd *cobra.Command, args []string) {
		inputName, _ := cmd.Flags().GetString("input")
		tileSize, _ := cmd.Flags().GetInt("tile-size")
		output, _ := cmd.Flags().GetString("output")

		mapPath, err := iohelpers.ResolveMapPath(inputName)
		if err != nil {
			fmt.Println("Map not found:", err)
			os.Exit(1)
		}

		fmt.Printf("Training tiles from %s (tile size: %d)...\n", mapPath, tileSize)
		if err := tiletrainer.Train(mapPath, tileSize, output); err != nil {
			fmt.Println("Training failed:", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Define command-line flags and bind to variables
	trainTilesCmd.Flags().StringVarP(&inputName, "input", "i", "", "Name of map to train on (without extension)")
	trainTilesCmd.Flags().IntVarP(&tileSize, "tile-size", "t", 64, "Size of each tile in pixels")
	trainTilesCmd.MarkFlagRequired("input")

	// Add the command to root
	rootCmd.AddCommand(trainTilesCmd)
}
