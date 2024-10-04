package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"sync"
	"os"
)

func FlipRight(img image.Image) image.Image {
	Y := img.Bounds().Max.Y
	X := img.Bounds().Max.X
	InverseX := Y
	reversed := image.NewNRGBA(image.Rect(0, 0, Y, X))

	for y := range Y {
		InverseY := X
		for x := range X {

			color := img.At(x, y)
			reversed.Set(InverseX, InverseY, reversed.ColorModel().Convert(color))
			InverseY--

		}

		InverseX--

	}
	return reversed

}

func FlipLeft(img image.Image) image.Image {
	Y := img.Bounds().Max.Y
	X := img.Bounds().Max.X
	InverseX := 0
	reversed := image.NewNRGBA(image.Rect(0, 0, Y, X))

	for y := range Y {
		InverseY := 0
		for x := range X {

			color := img.At(x, y)
			reversed.Set(InverseX, InverseY, reversed.ColorModel().Convert(color))
			InverseY++

		}

		InverseX++

	}
	return reversed

}

func FlipVertical(img image.Image) image.Image {
	Y := img.Bounds().Max.Y
	X := img.Bounds().Max.X
	InverseY := Y
	reversed := image.NewNRGBA(img.Bounds())

	for y := range Y {
		InverseX := X
		for x := range X {

			color := img.At(x, y)
			reversed.Set(InverseX, InverseY, reversed.ColorModel().Convert(color))
			InverseX--

		}

		InverseY--

	}
	return reversed

}

func OpenImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil

}

func SaveImage(path string, img image.Image) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := jpeg.Encode(file, img, nil); err != nil {
		return err
	}

	return nil

}

func main() {
	options := make(map[string]func(image.Image) image.Image)
	options["right"] = FlipRight
	options["left"] = FlipLeft
	options["vertical"] = FlipVertical

    operation := new(string)

	flag.StringVar(operation, "s","vertical", "Side to flip.Default=vertical")
    flag.Parse()
	args := flag.Args()

	var wg sync.WaitGroup

	for _, path := range args {
		wg.Add(1)
		go func() {

			defer wg.Done()
			img, err := OpenImage(path)
			if err != nil {
				fmt.Printf("Failed to open file %s: %s", path, err.Error())
				return

			}

			reversed := options[*operation](img)
			// newfile := fmt.Sprintf("%sreverted%d.jpg", path, idx)
			newfile := path

			if err := SaveImage(newfile, reversed); err != nil {

				fmt.Printf("failed to save file %s : %s ", newfile, err.Error())

			}

		}()
	}
	wg.Wait()

}
