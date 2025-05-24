package utils

import (
	"encoding/base64"
	"os"
	"strings"
)

func B64JsonToPng(b64Json string) error {
	b64data := b64Json

	// If the base64 string contains metadata (like "data:image/png;base64,"), remove it
	if i := strings.Index(b64data, ","); i != -1 {
		b64data = b64data[i+1:]
	}

	imageData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return err
	}

	err = os.WriteFile("output.png", imageData, 0644)
	if err != nil {
		return err
	}

	return nil
}
