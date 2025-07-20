package iohelpers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
)

var convertibleExts = map[string]func(io.Reader) (image.Image, error){
	".bmp":  bmp.Decode,
	".jpg":  jpeg.Decode,
	".jpeg": jpeg.Decode,
}

func ConvertToPNGIfNeeded(baseName string, folder string) (string, error) {
	// 1. Check if PNG already exists
	pngPath := filepath.Join(folder, baseName+".png")
	if _, err := os.Stat(pngPath); err == nil {
		return pngPath, nil
	}

	// 2. Try converting from known extensions
	for ext, decodeFn := range convertibleExts {
		for _, variant := range []string{ext, strings.ToUpper(ext)} {
			tryPath := filepath.Join(folder, baseName+variant)
			if _, err := os.Stat(tryPath); err == nil {
				fmt.Printf("Converting %s to PNG...\n", tryPath)
				in, _ := os.Open(tryPath)
				img, err := decodeFn(in)
				in.Close()
				if err != nil {
					return "", fmt.Errorf("decode failed: %v", err)
				}

				out, err := os.Create(pngPath)
				if err != nil {
					return "", fmt.Errorf("cannot create output: %v", err)
				}
				err = png.Encode(out, img)
				out.Close()
				if err != nil {
					return "", fmt.Errorf("PNG encode failed: %v", err)
				}

				// ✅ Conversion succeeded – delete the original
				err = os.Remove(tryPath)
				if err != nil {
					fmt.Printf("⚠️ Warning: PNG saved but failed to delete original: %v\n", err)
				} else {
					fmt.Printf("🗑️ Deleted original: %s\n", tryPath)
				}

				fmt.Printf("✅ Saved: %s\n", pngPath)
				return pngPath, nil
			}
		}
	}

	return "", fmt.Errorf("no image found for base name: %s", baseName)
}
