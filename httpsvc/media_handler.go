package httpsvc

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleUploadMedia(c echo.Context) error {
	fileInfo, err := c.FormFile("media")
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}
	src, err := fileInfo.Open()
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}
	defer src.Close()

	ext := filepath.Ext(fileInfo.Filename)
	fileName := utils.GenerateUniqueString() + ext
	dst := path.Join("media", fileName)

	err = s.uploader.Upload(dst, src)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusCreated, `{"status": "ok"}`)
}
