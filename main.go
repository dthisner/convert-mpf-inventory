package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	mpf "export-mountpf-inventory/MPF"
	"export-mountpf-inventory/models"
)

var smallsCollectionID = map[int]string{
	472543199509: "art-glass",
	286482694342: "containers",
	286482792646: "bells",
	471777870101: "copper",
	472025563413: "desktop-smalls",
	430688174357: "globes",
	472556405013: "trays",
	286460477638: "urns",
	286482596038: "gongs",
	286482727110: "hookahs",
	286482825414: "incense-burners",
	426051109141: "projectors",
	470517874965: "statues",
	286460608710: "teapot",
	287564398790: "tvs",
	286460641478: "animals",
	286482563270: "boxes",
	471777837333: "brass",
	286482759878: "candelabras",
	471777607957: "samovars",
	286460510406: "vases",
	285874553030: "cloisonne",
	471847371029: "inkwells",
	475696267541: "musical-instruments",
	471847305493: "pen-holders",
	286449172678: "sculptures",
	287093620934: "speakers",
	286482923718: "book-ends",
	470517907733: "busts",
	426050879765: "cameras",
	286482661574: "masks",
	471778787605: "navigation-equipment",
	286460575942: "planters",
	286482890950: "plates",
	287093653702: "turntables",
	286482858182: "bowls",
	286522441926: "clocks",
	472519999765: "medical",
	472735056149: "model-boats",
	286460543174: "figurines",
	286949703878: "radios",
	472520524053: "scientific",
	471707681045: "silver",
}

var lightningCollectionsID = map[int]string{
	471707615509: "bridge-lamp",
	286482759878: "candelabra",
	270290223302: "chandeliers",
	471707648277: "desk-lamp",
	470452470037: "floor-lamps",
	471707418901: "lamp-shades",
	270818279622: "neon-sign",
	270048460998: "scones",
	430756593941: "table-lamps",
	285402300614: "torchieres",
}

var (
	DUPLICATE_CHECK_JSON string     = "data/duplicateCheckCSV.json"
	DATA_MAP                        = make(map[string]bool)
	MUTEX                sync.Mutex // To ensure thread-safe updates
	TIME                 string
)

func main() {
	now := time.Now()
	TIME = now.Format("2006-01-02-15-04")

	// downloadRemainingImages()

	path := "export/CSV"
	listOfFiles, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var excelExports []models.Excel

	for _, file := range listOfFiles {
		filePath := filepath.Join(path, file.Name())

		excelExport, err := readCSVExport(filePath)
		if err != nil {
			log.Fatal(err)
		}

		excelExports = append(excelExports, excelExport...)
	}

	getCollections := false
	if getCollections {
		getCollectionIds()
		collections := mpf.GetCollections(smallsCollectionID, "smalls")
		// collections := mpf.GetCollectionsFromFolderWithJSON("./data/mpf")
		exportFromCollections(collections)
	}

	writeCSVFile(excelExports, "./export/csv/masterList.csv")

}

func cleanFilename(filename string) string {
	if idx := strings.Index(filename, "_"); idx != -1 {
		return filename[:idx]
	}
	return filename
}

func downloadRemainingImages() {
	var listOfDownloadedImages = make(map[string]bool)

	imagesPath := "/Users/sleipnir/Desktop/VPC Images"
	listOfImages, err := os.ReadDir(imagesPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, img := range listOfImages {
		fileName := cleanFilename(img.Name())

		if fileName == ".DS" {
			log.Printf("SKU Code: %s, skipping", fileName)
			continue
		}

		_, ok := listOfDownloadedImages[fileName]
		if !ok {
			log.Printf("Adding SKU: %s", fileName)
			listOfDownloadedImages[fileName] = true
		}
	}

	path := "export/CSV"
	listOfFiles, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var excelExports []models.Excel

	for _, file := range listOfFiles {
		filePath := filepath.Join(path, file.Name())

		excelExport, err := readCSVExport(filePath)
		if err != nil {
			log.Fatal(err)
		}

		excelExports = append(excelExports, excelExport...)
	}

	for _, row := range excelExports {
		_, ok := listOfDownloadedImages[row.Sku]

		if ok {
			log.Printf("SKU %s already exist", row.Sku)
			continue
		}

		for i, img := range row.Images {
			fileName := fmt.Sprintf("%s_0%d", row.Sku, i)

			err = downloadAndSaveImage(fileName, img.URL, "missing")
			if err != nil {
				log.Printf("Problem downloading SKU %s with URL: %s", row.Sku, fileName)
			}
		}
	}
}

func exportFromCollections(collections []models.CollectionData) {
	log.Printf("# of Collections %d", len(collections))

	for _, data := range collections {
		// "furniture-armoires-23-items"
		category := fmt.Sprintf("%s-%s-%d-items", data.Category, data.Name, data.TotalItems)
		log.Printf("category %s", category)

		// openJsonFileName := fmt.Sprintf("data/mpf/%s.json", category)
		exportCSVFileName := fmt.Sprintf("export/CSV/%s-%s.csv", TIME, category)
		exportJSONFileName := fmt.Sprintf("export/JSON/export-%s_%s.json", TIME, category)

		// MRF := openMRFJson(openJsonFileName)
		excelExport := generateExportData(data.MRP_DATA)
		DATA_MAP = openDuplicateCheckJson()
		missingSKU := 1

		for i, s := range excelExport {
			log.Printf(`Working with SKU: "%s" Item Number: %d out off: %d`, s.Sku, i, len(excelExport))
			excelExport[i].Completed = true
			excelExport[i].Duplicated = false
			excelExport[i].MissingSKU = false

			if len(s.Images) < 1 {
				log.Printf("No Images: %d", len(s.Images))
				excelExport[i].Completed = false
				excelExport[i].Error = "No Images to be located"
				continue
			}

			if s.Sku == "" {
				log.Printf("Missing SKU, here is image URL to find the item %s", s.Images[0].URL)
				s.Sku = fmt.Sprintf("%s%d", strings.ToUpper(data.Name), missingSKU)
				log.Printf("New SKU name is: %s", s.Sku)
				missingSKU++
			}

			if !isDuplicate(s.Sku) {
				updateDuplicateSkuMap(s.Sku)

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

func updateDuplicateSkuMap(key string) {
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
