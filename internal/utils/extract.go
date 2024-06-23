package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
	"github.com/flapperoo/beatmappacker/internal/logging"
)

func ExtractZip(dir string, url string, packNum uint16) (err error) {
	filePath, err := downloadArchive(dir, url, packNum)
	defer os.Remove(filePath)
	if err != nil {
		logging.BpLogger(logging.BpLog{
			Level:   "error",
			PackNum: packNum,
			Action:  "Downloading",
		})
		return err
	}

	logging.BpLogger(logging.BpLog{
		Level:   "info",
		PackNum: packNum,
		Action:  "Unpacking",
	})

	r, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer r.Close()
	// Unpack the file
	for _, f := range r.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		baseFileName := filepath.Base(f.Name)
		// Create the file
		dst, err := os.OpenFile(filepath.Join(dir, baseFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExtractSevenZip(dir string, url string, packNum uint16) (err error) {
	filePath, err := downloadArchive(dir, url, packNum)
	defer os.Remove(filePath)
	if err != nil {
		logging.BpLogger(logging.BpLog{
			Level:   "error",
			PackNum: packNum,
			Action:  "Downloading",
		})
		return err
	}

	logging.BpLogger(logging.BpLog{
		Level:   "info",
		PackNum: packNum,
		Action:  "Unpacking",
	})

	r, err := sevenzip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer r.Close()
	//Unpack the file
	for _, f := range r.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}
		// Extract base file name
		baseFileName := filepath.Base(f.Name)
		// Create the file
		dst, err := os.OpenFile(filepath.Join(dir, baseFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		defer dst.Close()
		// Open and Copy the file
		err = extractSevenZipFile(f, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractSevenZipFile(file *sevenzip.File, dst *os.File) (err error) {
	rc, err := file.Open()
	if err != nil {
		return err
	}

	defer rc.Close()
	// Copy the file
	_, err = io.Copy(dst, rc)
	if err != nil {
		return err
	}

	return nil
}
