package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
)

// 1 задание
func _first(ch <-chan int, w *sync.WaitGroup) {
	defer w.Done()
	for n := range ch {
		fmt.Printf("%d² = %d\n", n, n*n)
	}
}

// 2 задание
func filter(img draw.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := img.At(x, y).(color.RGBA)
			gray := uint8((int(oldColor.R) + int(oldColor.G) + int(oldColor.B)) / 3)
			newColor := color.RGBA{R: gray, G: gray, B: gray, A: oldColor.A}
			img.Set(x, y, newColor)
		}
	}
}

// 3 задание
func _filter(img draw.Image) {
	bounds := img.Bounds()
	var w sync.WaitGroup

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		w.Add(1)
		go func(y int) {
			defer w.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				oldColor := img.At(x, y).(color.RGBA)
				gray := uint8((int(oldColor.R) + int(oldColor.G) + int(oldColor.B)) / 3)
				newColor := color.RGBA{R: gray, G: gray, B: gray, A: oldColor.A}
				img.Set(x, y, newColor)
			}
		}(y)
	}
	w.Wait()
}

func main() {
	fmt.Println("Задание 1:")
	ch := make(chan int)
	var w sync.WaitGroup

	w.Add(1)
	go _first(ch, &w)

	for i := 1; i <= 9; i++ {
		ch <- i
	}
	close(ch)
	w.Wait()
	///////////////////////////////
	file, err := os.Open("dog.png")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	drawImg, ok := img.(draw.Image)
	if !ok {
		panic("failed")
	}

	start := time.Now()
	filter(drawImg)
	fmt.Println("Задание 2:", time.Since(start).String())

	outFile, err := os.Create("dog_2.png")
	if err != nil {
		panic(err)
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	err = png.Encode(outFile, drawImg)
	if err != nil {
		panic(err)
	}
	///////////////////////////////
	file, err = os.Open("dog.png")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	img, _, err = image.Decode(file)
	if err != nil {
		panic(err)
	}

	drawImg, ok = img.(draw.Image)
	if !ok {
		panic("failed")
	}

	start = time.Now()
	_filter(drawImg)
	fmt.Println("Задание 3:", time.Since(start).String())

	outFile, err = os.Create("dog_3.png")
	if err != nil {
		panic(err)
	}
	defer func(outFile *os.File) { _ = outFile.Close() }(outFile)

	err = png.Encode(outFile, drawImg)
	if err != nil {
		panic(err)
	}
}