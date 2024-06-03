package packer

import (
	"errors"
	"fmt"
	"os"

	"github.com/flapperoo/beatmappacker/internal/args"
	"github.com/flapperoo/beatmappacker/internal/logging"
	"github.com/flapperoo/beatmappacker/internal/utils"
)

func PackerProcess(a args.BpArgs) error {
	err := os.MkdirAll(a.PackPath, os.ModePerm)
	if err != nil {
		mkdirErr := errors.New("directory invalid")
		err = errors.Join(mkdirErr, err)
		return err
	}

	failedPacks := []uint64{}

	for i := a.FromPack; i <= a.ToPack; i++ {
		logging.BpLogger(logging.BpLog{
			Level:   "info",
			PackNum: i,
			Action:  "Processing",
		})

		var url, tempFile string
		var unpackFunc func(string, string) error

		switch {
		case i > 1299 || i == 124:
			if i > 1317 || i == 124 {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20osu%%21%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			} else {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			}
			tempFile = "temp.zip"
			unpackFunc = utils.UnzipZip
		default:
			url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20Beatmap%%20Pack%%20%%23%d.7z", i, i)
			tempFile = "temp.7z"
			unpackFunc = utils.UnzipSevenZip
		}
		// Download the file
		logging.BpLogger(logging.BpLog{
			Level:   "info",
			PackNum: i,
			Action:  "Downloading",
		})
		err := utils.Download(tempFile, url)
		if err != nil {
			logging.BpLogger(logging.BpLog{
				Level:   "error",
				PackNum: i,
				Action:  "Downloading",
				Err:     err,
			})
			logging.BpLogger(logging.BpLog{
				Level:   "infons",
				PackNum: i,
				Action:  "Skipping",
			})
			failedPacks = append(failedPacks, i)
			utils.CleanUp()
			continue
		}
		// Unpack the file
		logging.BpLogger(logging.BpLog{
			Level:   "info",
			PackNum: i,
			Action:  "Unpacking",
		})
		err = unpackFunc(tempFile, a.PackPath)
		if err != nil {
			logging.BpLogger(logging.BpLog{
				Level:   "error",
				PackNum: i,
				Action:  "Unpacking",
				Err:     err,
			})
			logging.BpLogger(logging.BpLog{
				Level:   "infons",
				PackNum: i,
				Action:  "Skipping",
			})
			failedPacks = append(failedPacks, i)
			utils.CleanUp()
			continue
		}

		logging.BpLogger(logging.BpLog{
			Level:   "infons",
			PackNum: i,
			Action:  "Finished",
		})
		utils.CleanUp()
	}
	// Print failed packs
	if len(failedPacks) > 0 {
		fmt.Println("[BeatmapPacker] Failed to process the following packs: ")
		for _, pack := range failedPacks {
			fmt.Printf("#%d ", pack)
		}
	}

	return nil
}
