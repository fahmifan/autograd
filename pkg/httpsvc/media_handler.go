package httpsvc

import (
	"net/http"

	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/mediastore"
	"github.com/fahmifan/autograd/pkg/core/mediastore/mediastore_cmd"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/labstack/echo/v4"
)

func (s *Server) handleSaveMedia(c echo.Context) error {
	authUser, ok := auth.GetUserFromCtx(c.Request().Context())
	if !ok {
		return responseError(c, ErrUnauthorized)
	}

	if !authUser.Role.Can(auth.CreateMedia) {
		return responseError(c, ErrUnauthorized)
	}

	fileInfo, err := c.FormFile("media")
	if err != nil {
		logs.ErrCtx(c.Request().Context(), err, "handleSaveMedia", "parse media")
		return responseError(c, err)
	}

	mediaType := c.FormValue("media_type")

	res, err := s.service.InternalSaveMultipart(c.Request().Context(), mediastore_cmd.InternalSaveMultipartRequest{
		FileInfo:  fileInfo,
		MediaType: mediastore.MediaFileType(mediaType),
	})
	if err != nil {
		logs.ErrCtx(c.Request().Context(), err, "handleSaveMedia", "InternalSaveMultipart")
		return responseError(c, err)
	}

	return c.JSON(http.StatusCreated, res)
}
