package packer

import (
	"errors"
	"fmt"
	"os"

	"github.com/flapperoo/beatmappacker/internal/args"
	"github.com/flapperoo/beatmappacker/internal/logging"
	"github.com/flapperoo/beatmappacker/internal/utils"
)

func Run(a args.BpArgs) error {
	err := os.MkdirAll(a.PackPath, os.ModePerm)
	if err != nil {
		mkdirErr := errors.New("directory invalid")
		err = errors.Join(mkdirErr, err)
		return err
	}

	failedPacks := []uint16{}

	for i := a.FromPack; i <= a.ToPack; i++ {
		logging.BpLogger(logging.BpLog{
			Level:   "info",
			PackNum: i,
			Action:  "Processing",
		})
		var url string
		switch {
		case i > 1299 || i == 124:
			if i > 1317 || i == 124 {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20osu%%21%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			} else {
				url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20Beatmap%%20Pack%%20%%23%d.zip", i, i)
			}
			err := utils.ExtractZip(a.PackPath, url, i)
			if err != nil {
				logging.BpLogger(logging.BpLog{
					Level:   "error",
					PackNum: i,
					Action:  "Processing",
					Err:     err,
				})
				logging.BpLogger(logging.BpLog{
					Level:   "infons",
					PackNum: i,
					Action:  "Skipping",
				})
				failedPacks = append(failedPacks, i)
				continue
			}
		default:
			url = fmt.Sprintf("https://packs.ppy.sh/S%d%%20-%%20Beatmap%%20Pack%%20%%23%d.7z", i, i)
			err := utils.ExtractSevenZip(a.PackPath, url, i)
			if err != nil {
				logging.BpLogger(logging.BpLog{
					Level:   "error",
					PackNum: i,
					Action:  "Processing",
					Err:     err,
				})
				logging.BpLogger(logging.BpLog{
					Level:   "infons",
					PackNum: i,
					Action:  "Skipping",
				})
				failedPacks = append(failedPacks, i)
				continue
			}
		}

		logging.BpLogger(logging.BpLog{
			Level:   "infons",
			PackNum: i,
			Action:  "Finished",
		})
	}
	// Print failed packs
	if len(failedPacks) > 0 {
		fmt.Print("[BeatmapPacker] Failed to process the following packs: ")
		for _, pack := range failedPacks {
			fmt.Printf("\n#%d ", pack)
		}
	}

	return nil
}
