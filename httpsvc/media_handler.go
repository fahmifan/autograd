package httpsvc

import (
	"net/http"

	"github.com/fahmifan/autograd/pkg/core/mediastore"
	"github.com/fahmifan/autograd/pkg/core/mediastore/mediastore_cmd"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleSaveMedia(c echo.Context) error {
	fileInfo, err := c.FormFile("media")
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	mediaType := c.FormValue("media_type")

	res, err := s.service.InternalSaveMultipart(c.Request().Context(), mediastore_cmd.InternalSaveMultipartRequest{
		FileInfo:  fileInfo,
		MediaType: mediastore.MediaFileType(mediaType),
	})
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusCreated, res)
}
