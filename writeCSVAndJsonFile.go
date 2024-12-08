package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

func writeCSVFile(excelExport []Excel, fileName string) {
	header := []string{"Sku", "Style", "Size", "Color", "Material", "Price", "Inventory", "Tags", "Image1", "Image2", "Image3", "Image4"}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for _, row := range excelExport {
		tags := strings.Join(row.Tags, ",")

		images := make([]string, 10) // Ensure 4 slots for images
		for i, img := range row.Images {
			if i >= 10 {
				break // Only include up to 4 images
			}
			images[i] = img.URL
		}

		// Prepare the row
		csvRow := []string{
			row.Sku,
			row.Descriptions.Style,
			row.Descriptions.Size,
			row.Descriptions.Color,
			row.Descriptions.Material,
			strconv.Itoa(row.Price),
			row.Inventory,
			tags,
		}
		csvRow = append(csvRow, images...) // Append images to the row

		if err := writer.Write(csvRow); err != nil {
			panic(err)
		}
	}
}

func writeJSONToFile(filename string, data interface{}) error {
	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create JSON encoder and write the data
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print with indentation
	return encoder.Encode(data)
}
