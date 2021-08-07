package model

import "io"

// ObjectStorer ..
type ObjectStorer interface {
	Store(dst string, r io.Reader) error
}
