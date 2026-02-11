# png2pic

Converts paletted PNG images to binary tileset format. Extracts tiles with configurable dimensions and spacing.

## Features

- Converts 256-color indexed PNG images to binary tile data format
- Configurable tile dimensions (e.g., 8x8, 16x16, 32x32)
- Support for tile spacing/padding
- Outputs compact binary format with header
- Debug mode for tile inspection

## Usage

### Command Line

```bash
png2pic [options]
```

### Options

- `-src <path>` - Path to the PNG image to convert (default: `tileset.png`)
- `-dst <path>` - Path to the output .pic file (default: `tileset.pic`)
- `-tile <format>` - Tile size in pixels (default: `16x16`)
- `-spacing <pixels>` - Spacing between tiles in pixels (default: `0`)
- `-debug` - Enable debug mode to show tile extraction details

### Examples

Convert a 16x16 tileset with no spacing:
```bash
png2pic -src mytiles.png -dst output.pic
```

Convert with 8x8 tiles and 1px spacing:
```bash
png2pic -src tileset.png -dst tiles.pic -tile 8x8 -spacing 1
```

Debug mode to inspect tile extraction:
```bash
png2pic -src tileset.png -debug
```

## Input Requirements

- Image must be a PNG file in indexed color mode (paletted)
- Maximum 256 colors in the palette
- Image dimensions should align with tile size + spacing

## Output Format

The output `.pic` file contains a binary header followed by tile data:

**Header (4 bytes):**
1. Tile width (uint8)
2. Tile height (uint8)
3. Bytes per tile (uint16, little-endian)
4. Total tiles count (uint8)

**Data:**
- Sequential tile data, row by row
- Each pixel is stored as a palette index (uint8)

## Installation

### Using Go

```bash
go install github.com/laghoule/png2pic@latest
```

### Using Docker

```bash
docker pull ghcr.io/laghoule/png2pic:latest
docker run -v $(pwd):/data ghcr.io/laghoule/png2pic -src /data/tileset.png -dst /data/output.pic
```

### From Source

```bash
git clone https://github.com/laghoule/png2pic.git
cd png2pic
go build -o png2pic main.go
```

## Use Cases

- Converting tileset graphics for retro game development
- Preparing tile data for embedded systems
- Creating compact binary tile formats for custom engines
- Processing sprite sheets with uniform tile dimensions
