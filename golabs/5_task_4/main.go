package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
)

func applyConvolutionKernel(src image.Image, kernel [][]float64, out *image.RGBA) {
	bounds := src.Bounds()
	kernelSize := len(kernel)
	kOffset := kernelSize / 2
	var wg sync.WaitGroup

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var rSum, gSum, bSum float64
				for ky := 0; ky < kernelSize; ky++ {
					for kx := 0; kx < kernelSize; kx++ {
						ix := x + kx - kOffset
						iy := y + ky - kOffset
						if ix >= bounds.Min.X && ix < bounds.Max.X && iy >= bounds.Min.Y && iy < bounds.Max.Y {
							pixel := src.At(ix, iy).(color.RGBA)
							weight := kernel[ky][kx]
							rSum += float64(pixel.R) * weight
							gSum += float64(pixel.G) * weight
							bSum += float64(pixel.B) * weight
						}
					}
				}
				r := uint8(clamp(rSum))
				g := uint8(clamp(gSum))
				b := uint8(clamp(bSum))
				out.Set(x, y, color.RGBA{R: r, G: g, B: b, A: src.At(x, y).(color.RGBA).A})
			}
		}(y)
	}
	wg.Wait()
}

func clamp(value float64) float64 {
	if value < 0 {
		return 0
	} else if value > 255 {
		return 255
	}
	return value
}

func normalizeKernel(kernel [][]float64) {
	sum := 0.0
	for _, row := range kernel {
		for _, value := range row {
			sum += value
		}
	}
	for y := range kernel {
		for x := range kernel[y] {
			kernel[y][x] /= sum
		}
	}
}

func main() {
	inputFile := "cat.png"
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// Преобразуем изображение в формат RGBA
	bounds := src.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, src, image.Point{}, draw.Src)

	// Ядро свёртки 7×7
	kernel := [][]float64{
		{1, 1, 2, 2, 2, 1, 1},
		{1, 2, 2, 4, 2, 2, 1},
		{2, 2, 4, 8, 4, 2, 2},
		{2, 4, 8, 16, 8, 4, 2},
		{2, 2, 4, 8, 4, 2, 2},
		{1, 2, 2, 4, 2, 2, 1},
		{1, 1, 2, 2, 2, 1, 1},
	}
	normalizeKernel(kernel)

	// Выходное изображение
	output := image.NewRGBA(bounds)

	// Cвёртку
	for i := 0; i < 5; i++ {
		applyConvolutionKernel(rgba, kernel, output)
	}

	outFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, output)
	if err != nil {
		panic(err)
	}

	fmt.Println("Фильтр успешно применён!")
}
