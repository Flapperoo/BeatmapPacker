package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/bodgit/sevenzip"
)

func main() {
	// Check if arguments are valid
	if len(os.Args) != 3 {
		log.Fatal("Invalid number of arguments, expected 2")
	}
	//Convert fromPack to int
	fromPack, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid argument %v", os.Args[1])
	}

	// Convert toPack to int
	toPack, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid argument %v", os.Args[2])
	}

	// Check if fromPack is greater than toPack
	if fromPack > toPack {
		log.Fatal("Invalid arguments, Initial pack number is greater than final pack number")
	}

	// Create Beatmap folder
	err = os.MkdirAll("BeatmapMegapack", os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating megapack directory: %v", err)
	}

	// Download & unpack beatmap packs
	for i := fromPack; i <= toPack; i++ {
		fmt.Printf("Downloading beatmap pack #%d\n", i)

		var url, tempFile string
		var unpackFunc func(string) error

		if i > 1299 || i == 124 { // TODO: Get edge case pack numbers from an external file instead of hardcoding them
			if i > 1317 || i == 124 {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20osu%%21%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			} else {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d20-%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			}
			tempFile = "temp.zip"
			unpackFunc = unzipZip
		} else {
			url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20Beatmap%%20Pack%%20%%23%d.7z", i, i)
			tempFile = "temp.7z"
			unpackFunc = unzipSevenZip
		}

		// Download the file
		err = download(tempFile, url)
		if err != nil {
			logAndCleanUp(fmt.Sprintf("Error downloading beatmap pack #%d", i), err)
			return
		}

		// Unpack the file
		err = unpackFunc(tempFile)
		if err != nil {
			logAndCleanUp(fmt.Sprintf("Error unpacking beatmap pack #%d", i), err)
			return
		}

		fmt.Printf("Beatmap Pack #%d unpacked\n", i)
	}

	//Clean up temp files
	logAndCleanUp("Beatmap unpacking finished, files are located in /BeatmapMegapack", nil)
}

func download(filepath string, url string) (err error) {
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

func unzipZip(filepath string) (err error) {
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
		dst, err := os.OpenFile(path.Join("BeatmapMegapack", baseFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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

func unzipSevenZip(filepath string) (err error) {
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
		dst, err := os.OpenFile(path.Join("BeatmapMegapack", baseFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
func logAndCleanUp(message string, err error) {
	fmt.Println(message)
	if err != nil {
		fmt.Println(err)
	}
	os.Remove("temp.zip")
	os.Remove("temp.7z")
}
