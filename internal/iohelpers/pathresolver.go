package iohelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveMapPath locates or converts a map into PNG format
func ResolveMapPath(name string) (string, error) {
	baseName := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))

	// 1. Check for existing .png in map-origins
	pngPath := filepath.Join("map_origins", baseName+".png")
	if _, err := os.Stat(pngPath); err == nil {
		return pngPath, nil
	}

	// 2. Try conversion
	convertedPath, err := ConvertToPNGIfNeeded(baseName, "map_origins")
	if err != nil {
		return "", fmt.Errorf("conversion failed for %s: %v", baseName, err)
	}

	return convertedPath, nil
}
