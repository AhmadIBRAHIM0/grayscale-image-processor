package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg" // Import necessary image format packages
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"
)

// Define a function to apply a grayscale filter to an image
func applyGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Iterate over the pixels of the input image and apply the grayscale formula to each pixel.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalPixel).(color.Gray)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

// Define a function to process an image
func processImage(inputPath, outputPath string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	// Open the input image file using os.Open.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
			return
		}
	}(inputFile)

	img, _, err := image.Decode(inputFile)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}

	// Apply the grayscale filter using the applyGrayscale function.
	grayImg := applyGrayscale(img)

	// Create the output file using os.Create.
	outputFile, err := os.Create(filepath.Join(outputPath, filepath.Base(inputPath)))
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
			return
		}
	}(outputFile)

	// Encode and save the grayscale image to the output file using the appropriate image format package.
	switch filepath.Ext(inputPath) {
	case ".png":
		err = png.Encode(outputFile, grayImg)
	case ".jpg":
		err = jpeg.Encode(outputFile, grayImg, nil)
	default:
		results <- fmt.Sprintf("Error when processing %s: Unsupported file format", inputPath)
		return
	}

	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}

	results <- fmt.Sprintf("Successfully processed %s, SUIII", inputPath)
}

func main() {
	inputPaths := []string{"assets/input/file1.jpg", "assets/input/file2.jpg", "assets/input/file3.jpg"}
	outputPath := "assets/output/"

	results := make(chan string)
	var wg sync.WaitGroup

	fmt.Println("Bom dia!, processing...")

	// Iterate over the results channel
	for _, inputPath := range inputPaths {
		wg.Add(1)
		go processImage(inputPath, outputPath, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Print the processing status for each image.
	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("Processing complete.")
	fmt.Println("Check the output folder for converted images.")
	fmt.Println("All rights reserved. © Ahmad IBRAHIM")
}
