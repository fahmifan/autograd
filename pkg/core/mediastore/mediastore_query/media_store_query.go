package mediastore_query

import (
	"context"
	"io"
	"path"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/mediastore"
	"github.com/google/uuid"
)

type MediaStoreQuery struct {
	*core.Ctx
}

type InternalFindMediaFileRequest struct {
	ID uuid.UUID
}

type InternalFindMediaFileResponse struct {
	BodyCloser io.ReadCloser
}

func (query *MediaStoreQuery) InternalFindMediaFile(ctx context.Context, req InternalFindMediaFileRequest) (InternalFindMediaFileResponse, error) {
	mediaFile, err := mediastore.MediaFileReader{}.FindByID(ctx, query.GormDB, req.ID.String())
	if err != nil {
		return InternalFindMediaFileResponse{}, err
	}

	srcPath := path.Join(query.RootDir, mediaFile.FilePath)
	readCloser, err := query.ObjectStorer.Seek(ctx, srcPath)
	if err != nil {
		return InternalFindMediaFileResponse{}, err
	}

	return InternalFindMediaFileResponse{
		BodyCloser: readCloser,
	}, nil
}
