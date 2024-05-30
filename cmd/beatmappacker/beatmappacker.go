package main

import (
	"fmt"
	"log"
	"os"

	"github.com/flapperoo/beatmappacker/internal/args"
	"github.com/flapperoo/beatmappacker/internal/packer"
)

func main() {
	// Assign arguments
	a, err := args.SetArgs(os.Args)
	if err != nil {
		log.Fatalf("Arguments invalid: %v", err)
	}

	err = packer.PackerProcess(a)
	if err != nil {
		log.Fatalf("Packer failed: %v", err)
	}

	fmt.Println("[BeatmapPacker] Beatmaps repacked!")
}
