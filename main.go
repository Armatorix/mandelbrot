package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

func main() {
	displayImage()
}

func displayImage() {
	const (
		// 1.5 x 1
		xmin, ymin, xmax, ymax = -2, -0.25, -1.25, 0.25
		width, height          = 30000, 20000
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup
	wg.Add(height)
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		go func(py int, y float64) {
			defer wg.Done()
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				img.Set(px, py, mandelbrot(z))
			}
		}(py, y)
	}
	f, err := os.Create("mandel.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

// https://github.com/budougumi0617/gopl/blob/master/ch03/ex05/mandelbrot.go
func mandelbrot(z complex128) color.Color {
	const iterations = 128
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			g := 255 - contrast*n
			r := g / 2
			b := contrast * n
			return color.RGBA{r, g, b, 255}
		}
	}
	return color.Black
}
