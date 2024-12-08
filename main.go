package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	openJsonFileName := "data/desks-149-items.json"
	exportCSVFileName := "export/CSV/desks-149-items.csv"
	exportJSONFileName := "export/JSON/desks-149-items.json"

	jsonFile, err := os.Open(openJsonFileName)
	if err != nil {
		log.Fatalf("issue opening file with err: %s", err)
	}

	byteValue, _ := io.ReadAll(jsonFile)

	var MRF MPF_EXPORT
	json.Unmarshal(byteValue, &MRF)
	excelExport := generateExportData(MRF)

	for i, s := range excelExport {
		excelExport[i].Completed = true

		for i, image := range s.Images {
			fileName := fmt.Sprintf("%s_0%d", s.Sku, i)
			err := downloadAndSaveImage(fileName, image.URL)
			if err != nil {
				log.Print(err.Error())
				s.Images[i].Error = err.Error()
				s.Images[i].Saved = false
				excelExport[i].Completed = false
			} else {
				s.Images[i].Saved = true
			}
		}
	}

	err = writeJSONToFile(exportJSONFileName, excelExport)
	if err != nil {
		fmt.Printf("Error writing JSON to file: %v\n", err)
	}

	writeCSVFile(excelExport, exportCSVFileName)
}

func generateExportData(MRF MPF_EXPORT) []Excel {
	var excelExport []Excel

	for _, s := range MRF.Data.Items {
		var excel Excel
		excel.Sku = s.Variants[0].Sku
		excel.Price = s.Variants[0].Price
		excel.Tags = removeIntFromTags(s.Tags)

		excel.Images = make([]ExcelImages, len(s.Images))
		for i, img := range s.Images {
			excel.Images[i].URL = strings.Replace(img.URL, "//", "https://", 1)
		}

		excel.Descriptions = extractDescriptions(s.Description)
		excel.Inventory = extractInventoryAmount(s.Title)

		excelExport = append(excelExport, excel)
	}

	return excelExport
}

func extractDescriptions(description string) Descriptions {
	var result Descriptions

	sizeRegex := regexp.MustCompile(`(?:Size|Dimensions):\s*([^<]+)`)
	materialRegex := regexp.MustCompile(`Material:\s*([^<]+)`)
	styleRegex := regexp.MustCompile(`(?:Style / Era|Style):\s*([^<]+)`)
	colorRegex := regexp.MustCompile(`Colour:\s*([^<]+)`)

	if match := sizeRegex.FindStringSubmatch(description); len(match) > 1 {
		result.Size = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(match[1]), "\"", ""))
	}

	if match := materialRegex.FindStringSubmatch(description); len(match) > 1 {
		result.Material = strings.ToUpper(strings.TrimSpace(match[1]))
	}

	if match := styleRegex.FindStringSubmatch(description); len(match) > 1 {
		result.Style = strings.ToUpper(strings.TrimSpace(match[1]))
	}

	if match := colorRegex.FindStringSubmatch(description); len(match) > 1 {
		result.Color = strings.ToUpper(strings.TrimSpace(match[1]))
	}

	return result
}

func extractInventoryAmount(title string) string {
	amountRegex := regexp.MustCompile(`x\d+`)

	match := amountRegex.FindString(title)
	if match != "" {
		return strings.ToUpper(match)
	}

	return "X1"
}

func removeIntFromTags(tags []string) []string {
	re := regexp.MustCompile(`^\d+$`)

	var result []string
	for _, item := range tags {
		if !re.MatchString(item) {
			result = append(result, strings.ToUpper(item))
		}
	}

	return result
}
