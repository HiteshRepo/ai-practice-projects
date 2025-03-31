package imageops

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func CreateMaskImage(originalImagePath string, x, y, width, height int) ([]byte, error) {
	originalImg, err := os.Open(originalImagePath)
	if err != nil {
		return nil, err
	}
	defer originalImg.Close()

	img, _, err := image.Decode(originalImg)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	mask := image.NewRGBA(bounds)

	draw.Draw(mask, bounds, &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Make the area to edit transparent
	for j := y; j < y+height && j < bounds.Max.Y; j++ {
		for i := x; i < x+width && i < bounds.Max.X; i++ {
			mask.Set(i, j, color.Transparent)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, mask); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
