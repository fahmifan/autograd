// Package fs provides implementation for uploading file
package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

type Config struct {
	_       string // enfore
	RootDir string
}

// LocalStorer ..
type LocalStorer struct {
	*Config
}

// NewLocalStorer ..
func NewLocalStorer(cfg *Config) *LocalStorer {
	return &LocalStorer{cfg}
}

// Store store file to local storage
func (f *LocalStorer) Store(dst, filename string, r io.Reader) error {
	dst = path.Join(f.RootDir, dst)
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("unable to create dst %s : %v", dst, err)
	}

	osFile, err := os.Create(path.Join(dst, filename))
	if err != nil {
		return fmt.Errorf("unable to create file %s : %v", dst, err)
	}

	_, err = io.Copy(osFile, r)
	if err != nil {
		return fmt.Errorf("unable to write file %s : %v", dst, err)
	}

	return nil
}

// Seek seek a file to local disk
func (f *LocalStorer) Seek(dst string) (r io.ReadCloser, err error) {
	return os.Open(path.Join(f.RootDir, dst))
}
