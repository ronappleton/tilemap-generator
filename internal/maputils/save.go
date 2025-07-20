package maputils

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
)

type TilesetEntry struct {
	ID        int       `json:"id"`
	File      string    `json:"file"`
	Hash      string    `json:"hash"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Adjacency Adjacency `json:"adjacency"`
}

type TilesetMetadata struct {
	TileSize int            `json:"tileSize"`
	Tiles    []TilesetEntry `json:"tiles"`
}

func SaveTileset(tiles []Tile, outputDir string, tileSize int) error {
	if err := os.MkdirAll(filepath.Join(outputDir, "tiles"), 0755); err != nil {
		return err
	}

	var entries []TilesetEntry

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

		entries = append(entries, TilesetEntry{
			ID:   tile.ID,
			File: filename,
			Hash: tile.Hash,
			X:    tile.X,
			Y:    tile.Y,
		})
	}

	metadata := TilesetMetadata{
		TileSize: tileSize,
		Tiles:    entries,
	}

	metaPath := filepath.Join(outputDir, "tileset.json")
	metaFile, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer metaFile.Close()

	return json.NewEncoder(metaFile).Encode(metadata)
}
