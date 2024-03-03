package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/bodgit/sevenzip"
)

func main() {
	if len(os.Args) == 3 {
		//Convert fromPack to int
		fromPack, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Invalid argument %v", os.Args[1])
			return
		}
		// Convert toPack to int
		toPack, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid argument %v", os.Args[2])
			return
		}
		// Check if fromPack is greater than toPack
		if fromPack > toPack {
			fmt.Println("Invalid arguments, Initial pack is greater than final pack")
			return
		}
		// Create Beatmap folder
		os.MkdirAll("BeatmapMegapack", os.ModePerm)
		// Download the packs
		for i := fromPack; i <= toPack; i++ {
			fmt.Println("Processing beatmap pack " + strconv.Itoa(i))
			url := ""
			if i > 1299 {
				if i > 1317 {
					url = "https://packs.ppy.sh/S" + strconv.Itoa(i) + "%20-%20osu%21%20Beatmap%20Pack%20%23" + strconv.Itoa(i) + ".zip"
				} else {
					url = "https://packs.ppy.sh/S" + strconv.Itoa(i) + "20-%20Beatmap%20Pack%20%23" + strconv.Itoa(i) + ".zip"
				}
				err := downloadFile("temp.zip", url)
				if err != nil {
					fmt.Println("Error downloading beatmap pack " + strconv.Itoa(i))
					fmt.Println("Does this link exist? " + url)
					fmt.Println(err)
					return
				}
				fmt.Println("Beatmap pack " + strconv.Itoa(i) + " downloaded")
				err = unpackBeatPackZip("temp.zip")
				if err != nil {
					fmt.Println("Error unpacking pack " + strconv.Itoa(i))
					fmt.Println(err)
					return
				}
			} else {
				url = "https://packs.ppy.sh/S" + strconv.Itoa(i) + "%20-%20Beatmap%20Pack%20%23" + strconv.Itoa(i) + ".7z"
				err := downloadFile("temp.7z", url)
				if err != nil {
					fmt.Println("Error downloading beatmap pack " + strconv.Itoa(i))
					fmt.Println("Does this link exist? " + url)
					return
				}
				fmt.Println("Beatmap pack " + strconv.Itoa(i) + " downloaded")
				err = unpackBeatPackSevenZip("temp.7z")
				if err != nil {
					fmt.Println("Error unpacking pack " + strconv.Itoa(i))
					fmt.Println(err)
					return
				}
			}
			fmt.Println("Beatmap Pack " + strconv.Itoa(i) + " unpacked")
		}
		fmt.Println("Cleaning up...")
		os.Remove("temp.zip")
		os.Remove("temp.7z")
		fmt.Println("Beatmap Megapack finished")
	} else {
		fmt.Println("Invalid arguments")
	}
}

func downloadFile(filepath string, url string) (err error) {
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

func unpackBeatPackZip(filepath string) (err error) {
	// Open the zip file
	file, err := zip.OpenReader(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Unpack the file
	for _, f := range file.File {
		// Create the file
		dst, err := os.OpenFile("BeatmapMegapack/"+f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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

func unpackBeatPackSevenZip(filepath string) (err error) {
	// Open the 7z file
	file, err := sevenzip.OpenReader(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	//Unpack the file
	for _, f := range file.File {
		// Create the file
		dst, err := os.OpenFile("BeatmapMegapack/"+f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Open the file
		fileShard, err := f.Open()
		if err != nil {
			return err
		}

		// Copy the file
		_, err = io.Copy(dst, fileShard)
		if err != nil {
			return err
		}
	}

	return nil
}
