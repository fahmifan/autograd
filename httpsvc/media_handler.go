package httpsvc

import (
	"fmt"
	"net/http"

	"github.com/fahmifan/autograd/config"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleUploadMedia(c echo.Context) error {
	fileInfo, err := c.FormFile("media")
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	fileName, err := s.mediaUsecase.Upload(fileInfo)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	publicURL := fmt.Sprintf("%s/media/%s", config.BaseURL(), fileName)
	return c.JSON(http.StatusCreated, map[string]string{"publicURL": publicURL})
}
