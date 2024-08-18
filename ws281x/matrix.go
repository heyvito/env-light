package ws281x

/*
#cgo CFLAGS: -std=c99 -I../../rpi_ws281x
#cgo LDFLAGS: -lws2811 -lm -L../../rpi_ws281x
#include <stdint.h>
#include <string.h>
#include <ws2811.h>
#include <wsgo.h>
*/
import "C"
import (
	"fmt"
	"image/color"
	"unsafe"
)

type Matrix struct {
	state      unsafe.Pointer
	color      color.RGBA
	brightness uint8
}

func (m *Matrix) Init(count, gpio int) error {
	var state unsafe.Pointer
	cErr := C.wsgo_init(C.int(count), C.int(gpio), &state)
	errNo := int(cErr)
	if errNo != 0 {
		return fmt.Errorf("wsgo_init errno: %d", errNo)
	}
	m.state = state
	return nil
}

func (m *Matrix) Finish() {
	C.wdgo_deinit(m.state)
}

func (m *Matrix) SetColor(r, g, b, brightness uint8) error {
	cErr := C.wsgo_set_color(m.state, r, g, b, brightness)
	errNo := int(cErr)
	if errNo != 0 {
		return fmt.Errorf("wsgo_set_color errno: %d", errNo)
	}

	return nil
}
