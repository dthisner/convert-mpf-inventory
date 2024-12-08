package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	DATA_MAP             = make(map[string]bool)
	MUTEX                sync.Mutex // To ensure thread-safe updates
	DUPLICATE_CHECK_JSON string     = "data/duplicateCheck.json"
)

func main() {
	fileName := "misc-4-items"
	openJsonFileName := fmt.Sprintf("data/mpf/%s.json", fileName)
	exportCSVFileName := fmt.Sprintf("export/CSV/%s.csv", fileName)
	exportJSONFileName := fmt.Sprintf("export/JSON/%s.json", fileName)

	MRF := openMRFJson(openJsonFileName)
	excelExport := generateExportData(MRF)

	DATA_MAP = openDuplicateCheckJson()

	for i, s := range excelExport {
		log.Printf(`Working with SKU: "%s"`, s.Sku)
		excelExport[i].Completed = true
		excelExport[i].Duplicated = false

		if !isDuplicate(s.Sku) {
			updateMap(s.Sku)

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

		} else {
			log.Printf(`SKU: "%s" is a duplicate`, s.Sku)
			excelExport[i].Duplicated = true
		}
	}

	writeToDuplicateCheckJson()
	err := writeJSONToFile(exportJSONFileName, excelExport)
	if err != nil {
		fmt.Printf("Error writing JSON to file: %v\n", err)
	}

	writeCSVFile(excelExport, exportCSVFileName)
}

func updateMap(key string) {
	MUTEX.Lock()
	defer MUTEX.Unlock()
	DATA_MAP[key] = true
	log.Printf("Updated key '%s' with value '%v'\n", key, true)
}

func isDuplicate(sku string) bool {
	if _, ok := DATA_MAP[sku]; ok {
		return true
	}

	return false
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
