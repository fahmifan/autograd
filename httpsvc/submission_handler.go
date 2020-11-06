package httpsvc

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

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
		return err
	}

	files := form.File["files"]

	for _, file := range files {

		src, err := file.Open()

		if err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

		defer src.Close()

		cwd, err := os.Getwd()

		if err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 16)
		randomString := generateRandomString(2)

		fileName := timestamp + "-" + randomString + "-" + file.Filename
		filePath := path.Join(cwd, "submission", fileName)
		dst, err := os.Create(filePath)

		if err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

	}

	return c.JSON(http.StatusOK, "success")
}

func generateRandomString(n int) string {
	letter := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)

	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}
