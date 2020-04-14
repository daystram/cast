package util

import (
	"errors"
	"fmt"
	"image"
	"os"

	"gitlab.com/daystram/cast/cast-be/config"

	"github.com/disintegration/imaging"
)

func NormalizeImage(root, filename string, width, height int) (err error) {
	var reader *os.File
	if reader, err = os.Open(fmt.Sprintf("%s/%s/%s.ori", config.AppConfig.UploadsDirectory, root, filename)); err != nil {
		return errors.New(fmt.Sprintf("[NormalizeProfile] failed to read original image. %+v", err))
	}
	original, _, err := image.Decode(reader)
	if err != nil {
		return
	}
	normalized := imaging.Fill(original, width, height, imaging.Center, imaging.Lanczos)
	if err = imaging.Save(normalized, fmt.Sprintf("%s/%s/%s.jpg", config.AppConfig.UploadsDirectory, root, filename)); err != nil {
		return errors.New(fmt.Sprintf("[NormalizeProfile] failed to normalize image. %+v", err))
	}
	reader.Close()
	_ = os.Remove(fmt.Sprintf("%s/%s/%s.ori", config.AppConfig.UploadsDirectory, root, filename))
	return
}
