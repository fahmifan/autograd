package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo"
)

func Test(context echo.Context) error {
	return context.String(http.StatusOK, fmt.Sprintf("this is a test"))
}

func Upload(c echo.Context) error {

	file, err := c.FormFile("file")

	if err != nil {
		return err
	}

	src, err := file.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	dst, err := os.Create(path.Join(cwd, "submission", file.Filename))

	fmt.Println(dst)

	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s is uploaded</p>", file.Filename))
}
