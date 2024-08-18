package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/heyvito/env-light/ws281x"
	"log"
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

	info := accessory.Info{Name: "Table Light"}
	ac := accessory.NewColoredLightbulb(info)
	ac.Lightbulb.Hue.OnValueRemoteUpdate(func(f float64) {
		fmt.Printf("Hue: %f\n", f)
	})
	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(f float64) {
		fmt.Printf("Saturation: %f\n", f)
	})
	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(i int) {
		fmt.Printf("Brightness: %d\n", i)
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
