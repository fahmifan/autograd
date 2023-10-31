package model

import (
	"context"
	"mime/multipart"
)

// MediaUsecase ..
type MediaUsecase interface {
	Upload(ctx context.Context, req *multipart.FileHeader) (name string, err error)
}
