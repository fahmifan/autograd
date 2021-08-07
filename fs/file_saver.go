// Package fs provides implementation for uploading file
package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

// LocalStorer ..
type LocalStorer struct{}

// NewLocalStorage ..
func NewLocalStorage() *LocalStorer {
	return &LocalStorer{}
}

// Store store file to local storage
func (f *LocalStorer) Store(dst string, r io.Reader) error {
	dir, _ := path.Split(dst)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("unable to create dir %s : %v", dir, err)
	}

	osFile, err := os.Create(dst)
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
	return os.Open(dst)
}
