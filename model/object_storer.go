package model

import "io"

// ObjectStorer ..
type ObjectStorer interface {
	Store(dst, filename string, r io.Reader) error
	Seek(dst string) (r io.ReadCloser, err error)
}
