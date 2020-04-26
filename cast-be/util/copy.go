package util

import (
	"io"
	"os"
)

// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang/21061062
func Copy(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
