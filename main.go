package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/heyvito/env-light/ws281x"
	"github.com/lucasb-eyer/go-colorful"
	"log"
)

/*
function hslToRgb(h, s, l) {
  var r, g, b;

  if (s == 0) {
    r = g = b = l; // achromatic
  } else {
    function hue2rgb(p, q, t) {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1/6) return p + (q - p) * 6 * t;
      if (t < 1/2) return q;
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
      return p;
    }

    var q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    var p = 2 * l - q;

    r = hue2rgb(p, q, h + 1/3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1/3);
  }

  return [ r * 255, g * 255, b * 255 ];
}
*/

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

func hslToRGB(h, s, l float64) (r, g, b float64) {
	if s == 0.0 {
		r, g, b = l, l, l
	} else {

		var q float64
		if l < 0.5 {
			q = l * (1.0 + s)
		} else {
			q = l + s - l*s
		}
		p := 2.0*l - q
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}

	r *= 255.0
	g *= 255.0
	b *= 255.0

	return
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mat := ws281x.Matrix{}
	err := mat.Init(78, 12)
	fatal(err)

	var H, S, L float64

	info := accessory.Info{Name: "Table Light"}
	ac := accessory.NewColoredLightbulb(info)

	updateColor := func() {
		fmt.Printf("Setting color HSL: %f, %f, %f\n", H, S, L)
		c := colorful.Hsv(H, S, L)
		//r, g, b := hslToRGB(H/360.0, S, L)
		r, g, b := c.RGB255()
		rr, gg, bb, ll := uint8(r), uint8(g), uint8(b), uint8(L*128.0)

		fmt.Printf("Setting color RGB: %d, %d, %d (L=%d)\n", rr, gg, bb, ll)
		if err := mat.SetColor(rr, gg, bb, ll); err != nil {
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
