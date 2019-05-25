package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var (
	width  = flag.Uint("w", 640, "width")
	height = flag.Uint("h", 0, "height")
	target = flag.String("target", "", "target folder")
)

func main() {
	flag.Parse()

	if *width == 0 && *height == 0 {
		log.Println("no width or height provided")
		return
	}

	dir := flag.Arg(0)
	if dir == "" {
		dir = "."
	}

	currentTarget := *target
	if currentTarget == "" {
		currentTarget = fmt.Sprintf("%v/%v_%v", dir, *width, *height)
		os.MkdirAll(currentTarget, os.ModePerm)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".png" && filepath.Ext(path) != ".jpg" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		m := resize.Resize(*width, *height, img, resize.Lanczos3)
		fn := filepath.Base(path)

		targetPath := fmt.Sprintf("%v/%v.jpg", currentTarget,
			strings.TrimSuffix(fn, filepath.Ext(fn)))
		out, err := os.Create(targetPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		jpeg.Encode(out, m, nil)

		return nil
	})
}
