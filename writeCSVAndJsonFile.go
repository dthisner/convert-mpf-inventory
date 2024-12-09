package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"export-mountpf-inventory/models"
)

func writeCSVFile(excelExport []models.Excel, fileName string) {
	header := []string{"Sku", "Style", "Size", "Color", "Material", "Price", "Inventory", "Tags", "Images"}

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
		var images []string

		for _, img := range row.Images {
			images = append(images, img.URL)
		}
		imageUrls := strings.Join(images, ",")

		csvRow := []string{
			row.Sku,
			row.Descriptions.Style,
			row.Descriptions.Size,
			row.Descriptions.Color,
			row.Descriptions.Material,
			strconv.Itoa(row.Price),
			row.Inventory,
			tags,
			imageUrls,
		}

		if err := writer.Write(csvRow); err != nil {
			panic(err)
		}
	}
}

func writeJSONToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print with indentation
	return encoder.Encode(data)
}

func openMRFJson(filename string) models.MPF_EXPORT {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("issue opening file with err: %s", err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var MRF models.MPF_EXPORT
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

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(m)
}
