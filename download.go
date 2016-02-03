package toolbelt

import (
	"io"
	"net/http"
	"os"
)

// Download file from URL and saves it in given path.
func Download(fromUrl string, toFilepath string) error {
	out, err := os.Create(toFilepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(fromUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
