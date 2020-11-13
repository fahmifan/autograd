package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleSubmission(c echo.Context) error {

	form, err := c.MultipartForm()

	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	files := form.File["files"]
	fileURLs := []string{}

	for _, file := range files {

		fileURL, err := s.submissionUsecase.Upload(file)

		if err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

		fileURLs = append(fileURLs, fileURL)
	}

	return c.JSON(http.StatusOK, map[string][]string{"file_urls": fileURLs})
}
