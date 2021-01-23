package model

import "mime/multipart"

// MediaUsecase ..
type MediaUsecase interface {
	Upload(fileInfo *multipart.FileHeader) (name string, err error)
}
