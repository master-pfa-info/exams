package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

const (
	output = "out.png"
	width  = 2048
	height = 1024
)

func main() {
	f, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	start := time.Now()
	img := create(width, height)
	delta := time.Since(start)
	fmt.Printf("time=%v\n", delta)

	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

// create fills one pixel at a time.
//
// time=??? <<< put the timing you find here.
func create(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			m.Set(i, j, pixel(i, j, width, height))
		}
	}
	return m
}

// create1 creates a Mandelbrot image.
//
// time=???
func create1(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	// ???
	return m
}

// create2 creates a Mandelbrot image.
//
// time=???
func create2(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	// ???
	return m
}

// create3 creates a Mandelbrot image.
//
// time=???
func create3(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	// ???
	return m
}

// pixel returns the color of a Mandelbrot fractal at the given point.
func pixel(i, j, width, height int) color.Color {
	const complexity = 1024

	xi := norm(i, width, -1.0, 2)
	yi := norm(j, height, -1, 1)

	const maxI = 1000
	x, y := 0., 0.

	for i := 0; (x*x+y*y < complexity) && i < maxI; i++ {
		x, y = x*x-y*y+xi, 2*x*y+yi
	}

	return color.Gray{uint8(x)}
}

func norm(x, total int, min, max float64) float64 {
	return (max-min)*float64(x)/float64(total) - max
}
