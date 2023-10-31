package mediastore_cmd

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/mediastore"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MediaStoreCmd struct {
	*core.Ctx
}

type InternalSaveMultipartRequest struct {
	FileInfo  *multipart.FileHeader
	MediaType mediastore.MediaFileType
}

type InternalSaveMultipartResponse struct {
	ID uuid.UUID `json:"id"`
}

func (m *MediaStoreCmd) InternalSaveMultipart(ctx context.Context, req InternalSaveMultipartRequest) (InternalSaveMultipartResponse, error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return InternalSaveMultipartResponse{}, core.ErrInternalServer
	}

	if !authUser.Role.Can(auth.CreateMedia) {
		return InternalSaveMultipartResponse{}, connect.NewError(connect.CodePermissionDenied, nil)
	}

	fileInfo := req.FileInfo

	src, err := fileInfo.Open()
	if err != nil {
		return InternalSaveMultipartResponse{}, err
	}
	defer src.Close()

	ext := filepath.Ext(fileInfo.Filename)
	fileName := utils.GenerateUniqueString() + ext
	dst := path.Join(m.RootFolder, fileName)

	err = m.ObjectStorer.Store(ctx, dst, src)
	if err != nil {
		logrus.Error(err)
		return InternalSaveMultipartResponse{}, err
	}

	publicURL := fmt.Sprintf("%s/%s", m.MediaServeBaseURL, fileName)

	mediaFile, err := mediastore.CreateMediaFile(mediastore.CreateMediaRequest{
		NewID:     uuid.New(),
		FileName:  fileInfo.Filename,
		FileType:  req.MediaType,
		Ext:       mediastore.Extension(ext),
		PublicURL: publicURL,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: mediastore.CreateMediaFile")
		return InternalSaveMultipartResponse{}, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err = mediastore.MediaFileWriter{}.Create(ctx, m.GormDB, &mediaFile)
	if err != nil {
		logs.ErrCtx(ctx, err, "MediaStoreCmd: InternalSaveMultipart: mediastore.MediaFileWriter{}.Create")
		return InternalSaveMultipartResponse{}, core.ErrInternalServer
	}

	return InternalSaveMultipartResponse{
		ID: mediaFile.ID,
	}, nil
}
