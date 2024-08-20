package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/heyvito/env-light/ws281x"
	"log"
)

type HSL struct {
	s, v, h float64
}

func (h HSL) RGB() (r, g, b, ww, cw float64) {
	s := h.s
	v := h.v
	if v == 0 {
		return
	}
	hue := h.h / 60.0
	const x = 255
	switch int(h.h / 60) {
	case 0:
		r, g, b = 1, hue, 0
	case 1:
		r, g, b = 2-hue, 1, 0
	case 2:
		b, r, g = hue-2, 0, 1
	case 3:
		r, g, b = 0, 4-hue, 1
	case 4:
		r, g, b = hue-4, 0, 1
	case 5:
		r, g, b = 1, 0, 6-hue
	}

	r = v * r
	g = v * g
	b = v * b
	cw = v * (1 - s) * 0.5

	return r * x, g * x, b * x, 0, cw * x
}

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
		c := HSL{H, S, L}
		r, g, b, _, cw := c.RGB()
		rr, gg, bb, ccw := uint8(r), uint8(g), uint8(b), uint8(cw)

		fmt.Printf("Setting color RGBL: %d, %d, %d, %d\n", rr, gg, bb, ccw)
		if err := mat.SetColor(rr, gg, bb, ccw); err != nil {
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
