package hash

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/corona10/goimagehash"
)

// PerceptualHash will return the perceptual for the image contained in 'r'.
func PerceptualHash(r io.Reader) (string, error) {

	im, _, err := image.Decode(r)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image, %w", err)
	}

	p_hash, err := goimagehash.PerceptionHash(im)

	if err != nil {
		return "", fmt.Errorf("Failed to derive perceptual hash, %w", err)
	}

	return p_hash.ToString(), nil
}
