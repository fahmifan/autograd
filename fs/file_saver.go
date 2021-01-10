// Package fs provides implementation for uploading file
package fs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// LocalStorage ..
type LocalStorage struct{}

// NewLocalStorage ..
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// Upload save file to local storage
func (f *LocalStorage) Upload(dst string, r io.Reader) error {
	dir, _ := path.Split(dst)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("unable to create dir %s : %v", dir, err)
	}

	bt, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("unable to read bytes: %v", err)
	}

	osFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("unable to create file %s : %v", dst, err)
	}

	_, err = osFile.Write(bt)
	if err != nil {
		return fmt.Errorf("unable to write file %s : %v", dst, err)
	}

	return nil
}
