package httpsvc

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/fahmifan/autograd/config"
	"github.com/fahmifan/autograd/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleUploadMedia(c echo.Context) error {
	mediaFile, err := c.FormFile("media")
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	file, err := mediaFile.Open()
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}
	defer file.Close()

	ext := filepath.Ext(mediaFile.Filename)
	fileName := utils.GenerateUniqueString() + ext

	err = s.objectStorer.Store("", fileName, file)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	publicURL := fmt.Sprintf("%s/api/v1/media/%s", config.APIBaseURL(), fileName)
	return c.JSON(http.StatusCreated, echo.Map{"publicURL": publicURL})
}

func (s *Server) handleGetMedia(c echo.Context) error {
	filename := c.Param("filename")
	media, err := s.objectStorer.Seek(filename)
	if err != nil {
		logrus.Error(err)
		return responseError(c, ErrNotFound)
	}
	defer media.Close()
	return c.Stream(http.StatusOK, "text/plain", media)
}
