package main

import (
	ws281x "github.com/mcuadros/go-rpi-ws281x"
	"image/color"
)

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config := &ws281x.DefaultConfig
	config.Brightness = 128
	config.Pin = 23
	width := 78
	height := 1

	c, err := ws281x.NewCanvas(width, height, config)
	fatal(err)

	defer c.Close()
	err = c.Initialize()
	fatal(err)

	bounds := c.Bounds()

	col := color.RGBA{255, 0, 0, 255}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c.Set(x, y, col)
		}
		err = c.Render()
		fatal(err)
	}
}
