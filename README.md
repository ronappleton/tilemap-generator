# Tilemap Generator

Tilemap Generator is a command line tool written in Go for analysing
bitmap maps and training reusable tilesets from them. It assists
creating procedural tilemaps by slicing a reference map into equal
sized tiles, deduplicating them and saving metadata for later
use. This README summarises how the application works and the
code organisation.

## Purpose and Overview

The tool's goal is to automate the creation of tile sets from an
existing image. The workflow is as follows:

1. A source map image is selected (BMP/JPEG/PNG supported).
2. The map is analysed to suggest a good tile size based on how many
   unique tiles each candidate size would yield.
3. The user chooses a tile size and the program slices the original
   image into tiles, deduplicates them and stores them in a new
   directory alongside `tileset.json` metadata.

## Command Line Interface

The CLI is implemented using Cobra. Available commands:

- `train-tiles` – analyse a map and generate a tileset.
- `list-maps`  – list images in `map_origins` ready for training.

`train-tiles` expects the map name (without extension) provided via
`--input` or `-i`. An optional `--diagnostic` flag saves a diagnostic
grid showing grouped tiles.

Typical usage:
```
# list available maps
./tilemap-generator list-maps

# train a tileset from example_map.png
./tilemap-generator train-tiles --input=example_map
```

The root command is simply `tilegen` as defined in `cmd/root.go`.

## Image Loading and Conversion

Images reside in the `map_origins/` folder. When a map name is given,
`iohelpers.ResolveMapPath` checks for an existing PNG. If only BMP or
JPEG versions are present it converts them using functions in
`image_converter.go` and deletes the originals. This ensures all maps
are processed in PNG format.

## Image Preprocessing

For analysis only (not slicing), images are cleaned via
`imagehelpers.PreprocessForTraining` which applies filters from the
`gift` package: grayscale conversion, a slight Gaussian blur and an
unsharp mask. Additional morphological operations remove noise and a
threshold step reduces anti‑aliasing. This preprocessing stabilises
hashing functions and tile comparison but the tiles written to disk are
cut from the untouched original image.

## Tile Size Analysis

`analyser.AnalyseTileSizesFuzzy` evaluates several candidate sizes. It
splits the preprocessed image into tiles and groups similar tiles using
fuzzy hashing. For each size the following metrics are collected:

- **TotalTiles** – number of tiles generated.
- **UniqueTiles** – count of distinct tiles after deduplication.
- **ReuseRatio** – proportion of tiles that are duplicates.

`PickSuggestedTileSize` returns the first size with a reuse ratio above
a threshold, suggesting a tile dimension that offers good reuse.

## Map Inspection

Prior to tile processing, `InspectMap` gathers statistics from the
cleaned image:

- Dimensions and colour model
- Number of unique colours and whether an alpha channel is used
- Average brightness and brightness range

These details are reported in the CLI before proceeding.

## Tileset Training

`tiletrainer.TrainFromImages` drives the tile extraction. The cleaned
image determines duplication groups; the original image provides the
pixels for saved tiles. SHA‑1 hashes are produced for deduplication via
`maputils.HashTile`. Tiles are stored as `tiles/tile_XXX.png` inside the
output directory.

If the `--diagnostic` flag is set, `SaveDiagnosticGrid` outputs a PNG
visualising tile groupings.

## Tileset JSON Output

`tileutils.SaveTilesetWithIndex` writes `tileset.json` containing:

- `tileSize` – the chosen tile dimension.
- `tiles` – array of entries with:
  - `id`       – unique tile ID
  - `file`     – relative path to tile image
  - `hash`     – SHA‑1 hash
  - `x`, `y`   – original grid coordinates
  - `adjacency` – neighbouring tile hashes (top, bottom, left, right)
- `mapping` – 2D array mapping positions in the original map to tile IDs.

This metadata allows later generation of new maps by referencing tiles
and understanding which tiles appeared adjacent in the source.

## Planned Extensions

The adjacency data in `tileset.json` can be used for procedural
generation. Improvements could include:

- Enhanced edge matching or auto‑tiling to respect terrain types.
- Symmetry detection to further reduce unique tile count.
- Rules describing which tiles may follow others for more controlled
  synthesis.

## File Structure

```
cmd/                 CLI commands
    root.go          Cobra root command
    list_maps.go     Lists available maps
    train_tiles.go   Main training workflow
internal/
    analyser/        Image inspection and tile size analysis
    imagehelpers/    Image loading and preprocessing filters
    iohelpers/       File format conversion and path resolution
    maputils/        Hashing, slicing and adjacency helpers
    tiletrainer/     High level training operations
    tileutils/       Tile extraction and saving logic
main.go              Entry point calling cmd.Execute()
```

## Dependencies

- [`github.com/spf13/cobra`](https://github.com/spf13/cobra)
  – command line framework.
- [`github.com/disintegration/gift`](https://github.com/disintegration/gift)
  – image filtering used for preprocessing.
- [`github.com/disintegration/imaging`](https://github.com/disintegration/imaging)
  – convenience functions for loading images.
- [`golang.org/x/image`](https://pkg.go.dev/golang.org/x/image)
  – BMP decoder and drawing utilities.
- [`github.com/corona10/goimagehash`](https://github.com/corona10/goimagehash)
  – perceptual hashing used by maputils (pHash).

The project targets Go 1.24 as specified in `go.mod`.
