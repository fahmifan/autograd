package grader

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// if outputPath is empty, the file will be download
// into temp with hashed sourceURL as file name
func download(outputPath, sourceURL string) error {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if outputPath == "" {
		hash := md5.New()
		hashed := hash.Sum([]byte(sourceURL))
		outputPath = fmt.Sprintf(path.Join(os.TempDir(), string(hashed)))
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
