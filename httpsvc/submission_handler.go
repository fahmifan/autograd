package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func responseError(c echo.Context, err error) error {
	switch err {
	case nil:
		return c.JSON(http.StatusOK, nil)
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) handleSubmission(c echo.Context) error {

	form, err := c.MultipartForm()

	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	files := form.File["files"]

	for _, file := range files {

		err := s.submissionUsecase.Upload(file)

		if err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

	}

	return c.JSON(http.StatusOK, "success")
}
