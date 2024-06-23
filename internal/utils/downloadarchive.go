package utils

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/flapperoo/beatmappacker/internal/logging"
)

func downloadArchive(dir string, genUrl string, packNum uint16) (filePath string, err error) {
	logging.BpLogger(logging.BpLog{
		Level:   "info",
		PackNum: packNum,
		Action:  "Downloading",
	})
	u, err := url.QueryUnescape(path.Base(genUrl))
	if err != nil {
		return "", err
	}

	fp := filepath.Join(dir, u)
	// Create the file
	f, err := os.Create(fp)
	if err != nil {
		return u, err
	}

	defer f.Close()
	// Get the data
	resp, err := http.Get(genUrl)
	if err != nil {
		return u, err
	}

	defer resp.Body.Close()
	// Write the body to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return u, err
	}

	return fp, nil
}
