package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

// copyAsset copies a named file from a source filesystem to a destination path.
func copyAsset(src fs.FS, name, dest string) error {
	in, err := src.Open(name)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	return errors.Join(err, out.Close())
}
