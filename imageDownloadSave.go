package main

import (
	"bufio"
	mpf "export-mountpf-inventory/MPF"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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

func downloadAndSaveImage(fileName, imageURL, category string) error {
	log.Printf("Downloading image %s", fileName)

	imgData, format, err := downloadImage(imageURL)
	if err != nil {
		return fmt.Errorf("issue: download image for %s with error: %s", fileName, err)
	}

	fileExtension := format
	if fileExtension == "jpeg" {
		fileExtension = "jpg" // standardize on ".jpg" for JPEG
	}

	folderPath := fmt.Sprintf("export/images/%s", category)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("issue creating directory: %v", err)
	}

	outputFilename := fmt.Sprintf("%s/%s.%s", folderPath, fileName, fileExtension)
	err = saveImage(imgData, outputFilename)
	if err != nil {
		return fmt.Errorf("issue: saving image for %s with error: %s", outputFilename, err)
	}

	log.Printf("Saved image to: %s", outputFilename)

	return nil
}

func getCollectionIds() {
	mpf.ReadEnvFile()

	urls := map[string]string{
		"/collections/animals":              "animals",
		"/collections/art-glass":            "art glass",
		"/collections/bell":                 "bells",
		"/collections/book-end":             "book ends",
		"/collections/bowl":                 "bowls",
		"/collections/box":                  "boxes",
		"/collections/brass?usf_take=84":    "brass",
		"/collections/busts":                "busts",
		"/collections/projectors":           "cameras",
		"/collections/candelabra":           "candelabras",
		"/collections/clocks":               "clocks",
		"/collections/cloisonne":            "cloisonne",
		"/collections/container":            "containers",
		"/collections/copper":               "copper",
		"/collections/desktop-smalls":       "desktop smalls",
		"/collections/figurine":             "figurines",
		"/collections/globers":              "globes",
		"/collections/gong":                 "gongs",
		"/collections/hookah":               "hookahs",
		"/collections/incense-burner":       "incense burners",
		"/collections/inkwells":             "inkwells",
		"/collections/mask":                 "masks",
		"/collections/medical-1":            "medical",
		"/collections/model-boats":          "model boats",
		"/collections/musical-instruments":  "musical instruments",
		"/collections/navigation-equipment": "navigation equipment",
		"/collections/pen-holders":          "pen holders",
		"/collections/planter":              "planters",
		"/collections/projectors-1":         "projectors",
		"/collections/plate":                "plates",
		"/collections/radios":               "radios",
		"/collections/samovar":              "samovars",
		"/collections/scientific":           "scientific",
		"/collections/sculptures":           "sculptures",
		"/collections/silver":               "silver",
		"/collections/speakers":             "speakers",
		"/collections/statue":               "statues",
		"/collections/teapot":               "teapot",
		"/collections/trays":                "trays",
		"/collections/turntables":           "turntables",
		"/collections/tvs":                  "tvs",
		"/collections/urn":                  "urns",
		"/collections/vase":                 "vases",
	}

	// Extracing the Collection ID when navigating to that collection
	BASE_URL := os.Getenv("BASE_URL")

	for u, name := range urls {
		url := fmt.Sprintf("%s%s", BASE_URL, u)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		imgData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(strings.NewReader(string(imgData)))
		for scanner.Scan() {
			// fmt.Println(scanner.Text())
			if strings.Contains(scanner.Text(), "_usfCollectionId") {
				reg := regexp.MustCompile(`\d+`)

				// Find the first occurrence of one or more digits
				match := reg.FindString(scanner.Text())
				fmt.Printf("%s:\"%s\",\n", match, name)
			}
		}
	}
}
