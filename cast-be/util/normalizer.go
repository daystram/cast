package util

import (
	"bytes"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

func NormalizeImage(input io.Reader, width, height int) (output bytes.Buffer, err error) {
	var original image.Image
	original, _, err = image.Decode(input)
	if err != nil {
		return
	}
	normalized := imaging.Fill(original, width, height, imaging.Center, imaging.Lanczos)
	err = imaging.Encode(&output, normalized, imaging.JPEG)
	return
}
