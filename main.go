package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/heyvito/env-light/ws281x"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math"
)

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mat := ws281x.Matrix{}
	err := mat.Init(78, 18)
	fatal(err)

	var H, S, L float64

	info := accessory.Info{Name: "Table Light"}
	ac := accessory.NewColoredLightbulb(info)

	updateColor := func() {
		fmt.Printf("Setting color HSL: %f, %f, %f\n", H, S, L)
		c := colorful.Hsl(H, S, L)
		r, g, b := c.RGB255()
		l := uint8(math.Ceil(L * 255.0))
		fmt.Printf("Setting color RGBL: %d, %d, %d, %d\n", r, g, b, l)
		if err := mat.SetColor(r, g, b, l); err != nil {
			fmt.Printf("Error updating color: %s\n", err)
		}
	}

	ac.Lightbulb.Hue.OnValueRemoteUpdate(func(f float64) {
		H = f // 0-360
		updateColor()
	})

	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(f float64) {
		S = f / 100.0
		updateColor()
	})

	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(i int) {
		L = float64(i) / 100.0
		updateColor()
	})

	ac.Lightbulb.On.OnValueRemoteUpdate(func(b bool) {
		if !b {
			if err := mat.SetColor(0, 0, 0, 0); err != nil {
				fmt.Printf("Error setting color: %s\n", err)
			}
		} else {
			updateColor()
		}
	})

	// configure the ip transport
	config := hc.Config{Pin: "00102003"}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	defer mat.Finish()
}
