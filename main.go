package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	tileFormat := flag.String("tile", "16x16", "Tile size in pixels")
	tileSpacing := flag.Uint("spacing", 0, "Spacing between tiles in pixels")
	srcName := flag.String("src", "tileset.png", "Path to the PNG image to convert")
	dstName := flag.String("dst", "tileset.pic", "Path to the output .pic file")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	fmt.Printf("png2pic version: %s, git commit: %s\n", version, gitCommit)

	tileWidth, tileHeight, err := extractTileFormat(*tileFormat)
	if err != nil {
		msg := fmt.Errorf("failed to extract tile format: %w", err)
		exitWithError(msg.Error())
	}

	pngFile, err := os.Open(*srcName)
	if err != nil {
		msg := fmt.Errorf("failed to open PNG file: %w", err)
		exitWithError(msg.Error())
	}
	defer pngFile.Close()

	datFile, err := os.Create(*dstName)
	if err != nil {
		msg := fmt.Errorf("failed to create .pic file: %w", err)
		exitWithError(msg.Error())
	}
	defer datFile.Close()

	img, err := png.Decode(pngFile)
	if err != nil {
		msg := fmt.Errorf("failed to decode PNG image: %w", err)
		exitWithError(msg.Error())
	}

	// g=Gt the indexed image
	pImg, indexed := img.(*image.Paletted)
	if !indexed {
		exitWithError("Image is not in indexed mode. Please use a PNG image with a 256-color palette.")
	}

	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

	// Must add spacing between tiles
	tilesByColumn := (imgWidth + int(*tileSpacing)) / (int(tileWidth) + int(*tileSpacing))
	tilesByRow := (imgHeight + int(*tileSpacing)) / (int(tileHeight) + int(*tileSpacing))

	if *debug {
		fmt.Println()
		fmt.Printf("TileSet dimension: %dx%d\n", imgWidth, imgHeight)
		fmt.Printf("Tiles by column: %d\n", tilesByColumn)
		fmt.Printf("Tiles by row: %d\n", tilesByRow)
		fmt.Printf("Number of tiles: %d\n", tilesByColumn*tilesByRow)
	}

	// Write header to pic file
	// TODO: create a picHeader struct
	binary.Write(datFile, binary.LittleEndian, tileWidth)                       // Tile width : 1 byte
	binary.Write(datFile, binary.LittleEndian, tileHeight)                      // Tile height : 1 byte
	binary.Write(datFile, binary.LittleEndian, 2*imgHeight)                     // Number of pixels in the image : 2 bytes
	binary.Write(datFile, binary.LittleEndian, uint8(tilesByColumn+tilesByRow)) // Number of tiles in the image : 1 byte

	for row := range tilesByRow {
		for col := range tilesByColumn {
			if *debug {
				fmt.Printf("\nColumn: %d, Row %d\n", col+1, row+1)
			}

			// Calculate tile position
			xStart := col * (int(tileWidth) + int(*tileSpacing))
			yStart := row * (int(tileHeight) + int(*tileSpacing))

			// Extract the tile
			for y := range int(tileHeight) {
				for x := range int(tileWidth) {
					var index uint8
					index = pImg.ColorIndexAt(xStart+x, yStart+y)
					binary.Write(datFile, binary.LittleEndian, index)
					if *debug {
						fmt.Printf("%03d ", index)
					}
				}
				if *debug {
					fmt.Println()
				}
			}
			// End row
		}
		// End col
	}

	fmt.Printf("Conversion complete! File generated: %s (%dx%d)\n", *dstName, tileWidth, tileHeight)
}

func extractTileFormat(format string) (uint8, uint8, error) {
	f := strings.ToLower(format)
	size := strings.Split(f, "x")

	if len(size) != 2 {
		return 0, 0, fmt.Errorf("invalid tile size format: %s", format)
	}

	width, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid tile width: %s", size[0])
	}

	height, err := strconv.Atoi(size[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid tile height: %s", size[1])
	}

	return uint8(width), uint8(height), nil
}

func exitWithError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
