package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	mpf "export-mountpf-inventory/MPF"
	"export-mountpf-inventory/models"
)

var (
	DUPLICATE_CHECK_JSON string     = "data/duplicateCheck.json"
	DATA_MAP                        = make(map[string]bool)
	MUTEX                sync.Mutex // To ensure thread-safe updates
)

func main() {
	// getCollectionIds()
	collections := mpf.GetCollections()

	for _, data := range collections {
		// "furniture-armoires-23-items"
		category := fmt.Sprintf("%s-%s-%d-items", data.Category, data.Name, data.TotalItems)
		// openJsonFileName := fmt.Sprintf("data/mpf/%s.json", category)
		exportCSVFileName := fmt.Sprintf("export/CSV/%s.csv", category)
		exportJSONFileName := fmt.Sprintf("export/JSON/export_%s.json", category)

		// MRF := openMRFJson(openJsonFileName)
		excelExport := generateExportData(data.MRP_DATA)

		DATA_MAP = openDuplicateCheckJson()

		missingSKU := 1

		for i, s := range excelExport {
			log.Printf(`Working with SKU: "%s" Item Number: %d`, s.Sku, i)
			excelExport[i].Completed = true
			excelExport[i].Duplicated = false
			excelExport[i].MissingSKU = false

			if s.Sku == "" {
				log.Printf("Missing SKU, here is image URL to find the item %s", s.Images[0].URL)
				s.Sku = fmt.Sprintf("%s%d", strings.ToUpper(data.Name), missingSKU)
				log.Printf("New SKU name is: %s")
				missingSKU++
			}

			if !isDuplicate(s.Sku) {
				updateMap(s.Sku)

				for i, image := range s.Images {
					fileName := fmt.Sprintf("%s_0%d", s.Sku, i)
					err := downloadAndSaveImage(fileName, image.URL, category)
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
			log.Printf("Error writing JSON to file: %v\n", err)
		}

		writeCSVFile(excelExport, exportCSVFileName)
	}
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

func generateExportData(MRF models.MPF_EXPORT) []models.Excel {
	var excelExport []models.Excel

	for _, s := range MRF.Data.Items {
		var excel models.Excel
		excel.Sku = s.Variants[0].Sku
		excel.Price = s.Variants[0].Price
		excel.Tags = removeIntFromTags(s.Tags)

		excel.Images = make([]models.ExcelImages, len(s.Images))
		for i, img := range s.Images {
			excel.Images[i].URL = strings.Replace(img.URL, "//", "https://", 1)
		}

		excel.Descriptions = extractDescriptions(s.Description)
		excel.Inventory = extractInventoryAmount(s.Title)

		excelExport = append(excelExport, excel)
	}

	return excelExport
}

func extractDescriptions(description string) models.Descriptions {
	var result models.Descriptions

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
