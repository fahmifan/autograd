package fs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// FileSaver ..
type FileSaver struct{}

// NewFileSaver ..
func NewFileSaver() *FileSaver {
	return &FileSaver{}
}

// Upload save file to local storage
func (f *FileSaver) Upload(dst string, r io.Reader) error {
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
