package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sync"
	"time"
)

var (
	width        int
	height       int
	maxIteration int
	fileName     string
)

// Command line arguments
func init() {
	flag.IntVar(&width, "width", 1920, "Image width.")
	flag.IntVar(&height, "height", 1080, "Image height.")
	flag.IntVar(&maxIteration, "epoch", 1000, "Amount of calculations per pixel.")
	flag.StringVar(&fileName, "fileName", "Mandelbrot", "File name to save.")
	flag.Parse()
}

func main() {
	// Create image
	img := image.NewGray16(image.Rectangle{Max: image.Point{X: width, Y: height}})
	// Create Wait Group for goroutines synchronization
	var wg sync.WaitGroup
	// Add delta to Wait Group
	wg.Add(width * height)

	fmt.Println("Start calculating")
	// Performance measurement
	start := time.Now()
	// For each pixel
	for px := 0; px < width; px++ {
		for py := 0; py < height; py++ {
			// Calculate Mandelbrot value in goroutine
			go mandelbrotXYGray16(px, py, img, &wg)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Finish performance measurement
	elapsed := time.Since(start)
	fmt.Printf("It took %f seconds", elapsed.Seconds())

	// Create file
	outputFile, err := os.Create(fileName + ".png")
	if err != nil {
		log.Panic(err)
	}

	// Save Image
	png.Encode(outputFile, img)
	outputFile.Close()
}

// Implementation of Mandelbrot algorithm
func mandelbrotXYGray16(px, py int, img *image.Gray16, wg *sync.WaitGroup) {
	// After algorithm executed tell WaitGroup that we finished
	defer wg.Done()

	// Scaling and vars initialization
	var x0 = scaleX(float64(px))
	var y0 = scaleY(float64(py))
	x := 0.0
	y := 0.0
	iteration := 0

	// Calculate the value maxIteration amount of times
	for x*x+y*y < 4.0 && iteration < maxIteration {
		var xtemp = x*x - y*y + x0
		y = 2*x*y + y0
		x = xtemp
		iteration = iteration + 1
	}

	// Transform calculation result to a greyscale color value
	grey := float64(iteration) * 65.535
	// Set image pixel to this color
	img.SetGray16(px, py, color.Gray16{Y: uint16(grey)})
}

func scaleX(x float64) float64 {
	//Mandelbrot X scale (-2.5, 1)
	f := 3.5*x/(float64(width)-1) - 2.5
	return f
}

func scaleY(y float64) float64 {
	//Mandelbrot Y scale (-1, 1)
	f := 2*y/(float64(height)-1) - 1
	return f
}
