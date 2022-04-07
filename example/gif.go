package main

import (
	"fmt"
	"image"
	"image/gif"
	"log"
	"os"
)

func outputGif(fileName string, images []*image.Paletted, delays []int) {
	err := os.MkdirAll("output", 0755)
	if err != nil {
		log.Fatal(err)
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
