package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"

	"github.com/DumDumGeniuss/ggol"
)

type DrawArea[T any] func(coord *ggol.Coordinate, area *T, unit int, image *image.Paletted, palette *[]color.Color)

func generateGif[T any](step int, unit int, duration int, fileName string, g *ggol.Game[T], drawArea DrawArea[T], palette []color.Color) {
	var images []*image.Paletted
	var delays []int
	var img *image.Paletted

	size := (*g).GetSize()

	for i := 0; i < step; i += 1 {
		img = image.NewPaletted(image.Rect(0, 0, size.Width*unit, size.Height*unit), palette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				area, _ := (*g).GetArea(coord)
				drawArea(coord, area, unit, img, &palette)
			}
		}
		images = append(images, img)
		delays = append(delays, duration)
		(*g).Iterate()
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
