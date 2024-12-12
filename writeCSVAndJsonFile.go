package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"export-mountpf-inventory/models"
)

func readCSVExport(filePath string) ([]models.Excel, error) {
	var MPFItems []models.Excel

	file, err := os.Open(filePath)
	if err != nil {
		return MPFItems, fmt.Errorf("Failed to open file: %s", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		return MPFItems, fmt.Errorf("Failed to Read file: %s", err)
	}

	// Skip the header if your file has one
	// If not, remove the next line
	records = records[1:]

	for _, record := range records {
		sku := record[0]

		if isDuplicate(sku) {
			log.Printf(`SKU: "%s" is a duplicate`, sku)
			continue
		}
		updateDuplicateSkuMap(sku)

		price, err := strconv.Atoi(record[5])
		if err != nil {
			log.Printf("Skipping record: %v (Failed to convert age to int: %s)\n", record, err)
			continue
		}

		tags := strings.Split(record[7], ",")
		splitImages := strings.Split(record[8], ",")

		var images []models.ExcelImages
		for _, img := range splitImages {
			var image models.ExcelImages

			image.URL = img
			images = append(images, image)
		}

		desc := models.Descriptions{
			Style:    record[1],
			Size:     record[2],
			Color:    record[3],
			Material: record[4],
		}

		item := models.Excel{
			Sku:          sku,
			Price:        price,
			Inventory:    record[6],
			Tags:         tags,
			Images:       images,
			Descriptions: desc,
		}

		MPFItems = append(MPFItems, item)
	}

	writeToDuplicateCheckJson()

	return MPFItems, nil
}

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
