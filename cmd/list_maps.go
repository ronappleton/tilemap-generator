package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tilemap-generator/internal/iohelpers"

	"github.com/spf13/cobra"
)

var listMapsCmd = &cobra.Command{
	Use:   "list-maps",
	Short: "List and select a map from the map_origins/ folder (converts if needed)",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := os.ReadDir("map_origins")
		if err != nil {
			fmt.Println("Failed to read map_origins/:", err)
			return
		}

		var files []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			lower := strings.ToLower(name)
			if strings.HasSuffix(lower, ".png") ||
				strings.HasSuffix(lower, ".bmp") ||
				strings.HasSuffix(lower, ".jpg") ||
				strings.HasSuffix(lower, ".jpeg") {
				files = append(files, name)
			}
		}

		if len(files) == 0 {
			fmt.Println("No valid image maps found in map-origins/")
			return
		}

		fmt.Println("Available maps:")
		for i, f := range files {
			fmt.Printf(" [%d] %s\n", i+1, f)
		}

		fmt.Print("\nSelect a map number: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var index int
		_, err = fmt.Sscanf(input, "%d", &index)
		if err != nil || index < 1 || index > len(files) {
			fmt.Println("Invalid selection.")
			return
		}

		selected := files[index-1]
		name := strings.TrimSuffix(selected, filepath.Ext(selected))

		// Force convert to PNG if needed
		resolvedPath, err := iohelpers.ResolveMapPath(name)
		if err != nil {
			fmt.Println("Conversion failed:", err)
			return
		}

		fmt.Printf("\n✅ You selected: %s\n", filepath.Base(resolvedPath))
		fmt.Printf("👉 Use it with:\n  tilemap-generator train-tiles --input=%s\n", name)
	},
}

func init() {
	rootCmd.AddCommand(listMapsCmd)
}
