package main

import (
	"image/png"
	"minilight/radiance"
	"os"
)

func main() {
	ri := radiance.NewRadianceImage(500, 500)
	for x := 0; x < 500; x++ {
		for y := 0; y < 500; y++ {
			ri.Add(x, y, radiance.Radiance{1, 0, 0})
		}
	}
	im := ri.ToRGBA64(1)

	w, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(w, im)
}
