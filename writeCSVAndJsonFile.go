package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
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

func openMRFJson(filename string) MPF_EXPORT {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("issue opening file with err: %s", err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var MRF MPF_EXPORT
	json.Unmarshal(byteValue, &MRF)

	return MRF
}

func openDuplicateCheckJson() map[string]bool {
	jsonFile, err := os.Open(DUPLICATE_CHECK_JSON)
	duplicateCheck := make(map[string]bool)

	if err != nil {
		return duplicateCheck
	}

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &duplicateCheck)

	return duplicateCheck
}

func writeToDuplicateCheckJson() {
	err := writeMapToFile(DUPLICATE_CHECK_JSON, DATA_MAP)
	if err != nil {
		log.Printf("Error writing updated map to file: %v\n", err)
	}
}

func writeMapToFile(filename string, m map[string]bool) error {
	MUTEX.Lock()
	defer MUTEX.Unlock()

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serialize the map to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(m)
}
