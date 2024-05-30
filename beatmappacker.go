package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	// Check if arguments are valid
	if len(os.Args) != 3 {
		log.Fatal("Invalid number of arguments, expected 2")
	}

	fromPack, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("Invalid argument %v", os.Args[1])
	}

	toPack, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("Invalid argument %v", os.Args[2])
	}

	if fromPack > toPack {
		log.Fatal("Invalid arguments, Initial pack number is greater than final pack number")
	}

	// Create Beatmap folder
	err = os.MkdirAll("BeatmapMegapack", os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating megapack directory: %v", err)
	}

	// Track failed packs
	failedPacks := []uint64{}

	// BeatmapPacker Main
	for i := fromPack; i <= toPack; i++ {
		bpLogger(logOptions{
			Level:   "info",
			PackNum: i,
			Action:  "Processing",
		})

		var url, tempFile string
		var unpackFunc func(string) error

		switch {
		case i > 1299 || i == 124:
			if i > 1317 || i == 124 {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20osu%%21%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			} else {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d20-%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			}
			tempFile = "temp.zip"
			unpackFunc = unzipZip
		default:
			url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20Beatmap%%20Pack%%20%%23%d.7z", i, i)
			tempFile = "temp.7z"
			unpackFunc = unzipSevenZip
		}

		// Download the file
		bpLogger(logOptions{
			Level:   "info",
			PackNum: i,
			Action:  "Downloading",
		})
		err = download(tempFile, url)
		if err != nil {
			bpLogger(logOptions{
				Level:   "error",
				PackNum: i,
				Action:  "Downloading",
				Err:     err,
			})
			bpLogger(logOptions{
				Level:   "infons",
				PackNum: i,
				Action:  "Skipping",
			})
			failedPacks = append(failedPacks, i)
			cleanUp()
			continue
		}

		// Unpack the file
		bpLogger(logOptions{
			Level:   "info",
			PackNum: i,
			Action:  "Unpacking",
		})
		err = unpackFunc(tempFile)
		if err != nil {
			bpLogger(logOptions{
				Level:   "error",
				PackNum: i,
				Action:  "Unpacking",
				Err:     err,
			})
			bpLogger(logOptions{
				Level:   "infons",
				PackNum: i,
				Action:  "Skipping",
			})
			failedPacks = append(failedPacks, i)
			cleanUp()
			continue
		}

		bpLogger(logOptions{
			Level:   "infons",
			PackNum: i,
			Action:  "Finished",
		})
		cleanUp()
	}

	// Print failed packs
	if len(failedPacks) > 0 {
		fmt.Println("[BeatmapPacker] Failed to process the following packs: ")
		for _, pack := range failedPacks {
			fmt.Printf("#%d ", pack)
		}
	}

	fmt.Println("[BeatmapPacker] Beatmaps repacked!")
}
