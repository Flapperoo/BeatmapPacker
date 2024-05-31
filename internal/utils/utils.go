package utils

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/bodgit/sevenzip"
)

func Download(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func UnzipZip(filepath string, packPath string) (err error) {
	// Open the zip file
	file, err := zip.OpenReader(filepath)
	if err != nil {
		return err
	}

	defer file.Close()
	// Unpack the file
	for _, f := range file.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}
		// Extract base file name
		baseFileName := path.Base(f.Name)
		// Create the file
		dst, err := os.OpenFile(path.Join(packPath, baseFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		defer dst.Close()
		// Open the file
		fileShard, err := f.Open()
		if err != nil {
			return err
		}

		defer fileShard.Close()
		// Copy the file
		_, err = io.Copy(dst, fileShard)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnzipSevenZip(filepath string, packPath string) (err error) {
	// Open the 7z file
	file, err := sevenzip.OpenReader(filepath)
	if err != nil {
		return err
	}

	defer file.Close()
	//Unpack the file
	for _, f := range file.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}
		// Extract base file name
		baseFileName := path.Base(f.Name)
		// Create the file
		dst, err := os.OpenFile(path.Join(packPath, baseFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		defer dst.Close()
		// Open and Copy the file
		err = processSevenZipFile(f, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func CleanUp() {
	os.Remove("temp.zip")
	os.Remove("temp.7z")
}

func processSevenZipFile(file *sevenzip.File, dst *os.File) (err error) {
	fileShard, err := file.Open()
	if err != nil {
		return err
	}

	defer fileShard.Close()
	// Copy the file
	_, err = io.Copy(dst, fileShard)
	if err != nil {
		return err
	}

	return nil
}
