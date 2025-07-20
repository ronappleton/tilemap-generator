package tileutils

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"tilemap-generator/internal/maputils"
)

// SaveTilesetWithIndex saves unique tiles to disk and writes a metadata file
// containing the mapping of tile positions to tile IDs.
func SaveTilesetWithIndex(tiles []maputils.Tile, mapping [][]int, outputDir string, tileSize int) error {
	if err := os.MkdirAll(filepath.Join(outputDir, "tiles"), 0755); err != nil {
		return err
	}

	adj := maputils.BuildAdjacency(tiles, mapping)

	var entries []maputils.TilesetEntry
	for _, tile := range tiles {
		filename := fmt.Sprintf("tiles/tile_%03d.png", tile.ID)
		fullPath := filepath.Join(outputDir, filename)

		outFile, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		if err := png.Encode(outFile, tile.Image); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()

		entries = append(entries, maputils.TilesetEntry{
			ID:        tile.ID,
			File:      filename,
			Hash:      tile.Hash,
			X:         tile.X,
			Y:         tile.Y,
			Adjacency: adj[tile.ID],
		})
	}

	meta := struct {
		TileSize int                     `json:"tileSize"`
		Tiles    []maputils.TilesetEntry `json:"tiles"`
		Mapping  [][]int                 `json:"mapping"`
	}{
		TileSize: tileSize,
		Tiles:    entries,
		Mapping:  mapping,
	}

	metaPath := filepath.Join(outputDir, "tileset.json")
	f, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(meta)
}
