package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"

	"github.com/DumDumGeniuss/ggol"
)

const (
	BlackColorIndex = iota
	WhiteColorIndex
	RedColorIndex
	BlueColorIndex
	CyanColorIndex
	GreenColorIndex
	YellowColorIndex
	OrangeColorIndex
	PurpleColorIndex
	GoldColorIndex
)

type DrawCell[T any] func(coord *ggol.Coordinate, cell *T, unit int, image *image.Paletted, palette *[]color.Color)

func generateGif[T any](duration int, fileName string, g *ggol.Game[T], drawCell DrawCell[T]) {
	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xe5, 0x39, 0x35, 0xff},
		color.RGBA{0x1e, 0x88, 0xe5, 0xff},
		color.RGBA{0x00, 0xac, 0xc1, 0xff},
		color.RGBA{0x43, 0xa0, 0x47, 0xff},
		color.RGBA{0xfd, 0xd8, 0x35, 0xff},
		color.RGBA{0xfb, 0x8c, 0x00, 0xff},
		color.RGBA{0x8e, 0x24, 0xaa, 0xff},
		color.RGBA{0xff, 0xd7, 0x00, 0xff},
	}
	var images []*image.Paletted
	var delays []int
	var img *image.Paletted
	var unit int = 10

	size := (*g).GetSize()

	for step := 0; step < 100; step += 1 {
		img = image.NewPaletted(image.Rect(0, 0, size.Width*unit, size.Height*unit), palette)
		for x := 0; x < size.Width; x += 1 {
			for y := 0; y < size.Height; y += 1 {
				coord := &ggol.Coordinate{X: x, Y: y}
				cell, _ := (*g).GetCell(coord)
				drawCell(coord, cell, unit, img, &palette)
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
