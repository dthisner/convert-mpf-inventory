package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // Import JPEG decoder
	_ "image/png"  // Import PNG decoder
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "golang.org/x/image/webp" // Import WebP decoder
)

func downloadImage(url string) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Use image.Decode to identify the format
	_, format, err := image.Decode(strings.NewReader(string(imgData)))
	if err != nil {
		return nil, "", err
	}

	// Ensure the format is not empty
	if format == "" {
		return nil, "", fmt.Errorf("could not identify image format")
	}

	return imgData, format, nil
}

func saveImage(imgData []byte, filename string) error {
	// Create the file to save the image
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write the image data to the file
	_, err = outFile.Write(imgData)
	if err != nil {
		return err
	}
	return nil
}

func downloadAndSaveImage(fileName, url string) error {
	imageURL := url
	log.Printf("Downloading image %s", fileName)

	imgData, format, err := downloadImage(imageURL)
	if err != nil {
		return fmt.Errorf("issue: download image for %s with error: %s", fileName, err)
	}

	fileExtension := format
	if fileExtension == "jpeg" {
		fileExtension = "jpg" // standardize on ".jpg" for JPEG
	}

	outputFilename := fmt.Sprintf("export/images/%s.%s", fileName, fileExtension)
	err = saveImage(imgData, outputFilename)
	if err != nil {
		return fmt.Errorf("issue: saving image for %s with error: %s", outputFilename, err)
	}

	return nil
}
