package logging

import "github.com/fatih/color"

type BpLog struct {
	Level   string
	PackNum uint16
	Action  string
	Err     error
}

func BpLogger(log BpLog) {
	w := color.New(color.FgWhite)
	r := color.New(color.FgRed)
	switch log.Level {
	case "info":
		w.Printf("\033[2K\r[BeatmapPacker] %v Beatmap Pack #%d", log.Action, log.PackNum)
	case "infons":
		w.Printf("\033[2K\r[BeatmapPacker] %v Beatmap Pack #%d\n", log.Action, log.PackNum)
	case "error":
		w.Printf("\n[BeatmapPacker] ")
		r.Print("(Error) ")
		w.Printf("%v Beatmap Pack #%d:\n %v\n", log.Action, log.PackNum, log.Err)
	}
}
