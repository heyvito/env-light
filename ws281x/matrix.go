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
	state      any
	color      color.RGBA
	brightness uint8
}

func (m *Matrix) Init(count, gpio int) error {
	err := C.wsgo_init(unsafe.Pointer(&m.state), C.int(count), C.int(gpio))
	fmt.Printf("Error is %#v\n", err)

	return nil
}
