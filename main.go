package main

import "github.com/heyvito/env-light/ws281x"

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mat := ws281x.Matrix{}
	err := mat.Init(78, 18)
	fatal(err)

	err = mat.SetColor(144, 0, 217, 128)
	fatal(err)

	defer mat.Finish()
}
