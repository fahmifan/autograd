// Package iface provides general interfaces for package uscease.
package iface

import "mime/multipart"

type (
	// MediaUsecase ..
	MediaUsecase interface {
		Upload(fileInfo *multipart.FileHeader) (pubURL string, err error)
	}
)
