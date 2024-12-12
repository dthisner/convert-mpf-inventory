package mpf

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"export-mountpf-inventory/models"

	"github.com/joho/godotenv"
)

var (
	API_KEY string
	API_URL string
)

func ReadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	API_KEY = os.Getenv("API_KEY")
	API_URL = os.Getenv("API_URL")
}

func getCollectionData(id, start, end int) (models.MPF_EXPORT, error) {
	var MPF models.MPF_EXPORT

	ReadEnvFile()
	log.Printf("Getting for Collection ID: %d from %d to %d", id, start, end)
	url := fmt.Sprintf("%ssearch?q=&apiKey=%s&locale=en&collection=%d&skip=%d&take=%d&sort=title", API_URL, API_KEY, id, start, end)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return MPF, fmt.Errorf("error making request: %v", err)

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MPF, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&MPF); err != nil {
		return MPF, fmt.Errorf("error decoding JSON: %v", err)
	}

	return MPF, nil
}

func GetCollections(collectionName map[int]string, category string) []models.CollectionData {
	var collectionData []models.CollectionData

	for id, name := range collectionName {
		var collection models.CollectionData
		start := 0
		end := 200

		MPF, err := getCollectionData(id, start, end)
		if err != nil {
			log.Print(err)
		}

		collection.MRP_DATA = MPF
		collection.ID = id
		collection.Category = category
		collection.Name = name
		collection.TotalItems = MPF.Data.Total

		log.Printf("MPF.Data.Total: %d End: %d", MPF.Data.Total, end)
		if MPF.Data.Total > end {
			for {
				start = start + 200
				end = end + 200

				log.Printf("MPF.Data.Total: %d End: %d id: %d", MPF.Data.Total, end, id)
				MPF, err := getCollectionData(id, start, end)
				if err != nil {
					log.Print(err)
				}

				collection.MRP_DATA.Data.Items = append(collection.MRP_DATA.Data.Items, MPF.Data.Items...)

				if MPF.Data.Total < end {
					break
				}
			}
		}

		collectionData = append(collectionData, collection)
	}

	fileName := fmt.Sprintf("./data/%s.json", category)
	writeJSONToFile(fileName, collectionData)

	return collectionData
}

func GetCollectionsFromFolderWithJSON(path string) []models.CollectionData {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var collectionData []models.CollectionData

	for _, file := range files {
		var collection models.CollectionData
		var MPF models.MPF_EXPORT

		log.Printf("file: \"%s\"", file)

		filePath := filepath.Join(path, file.Name())
		fmt.Printf("Processing file: %s\n", filePath)

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}

		if err := json.Unmarshal([]byte(content), &MPF); err != nil {
			log.Printf("Failed to unmarshal file %s: %v", filePath, err)
			continue
		}

		result := extractCategoryAndCollection(file.Name())

		collection.MRP_DATA = MPF
		collection.Category = result["category"]
		collection.Name = result["collection"]
		collection.TotalItems = MPF.Data.Total

		log.Printf("category: \"%s\" collection: \"%s\" Total: \"%d\"", collection.Category, collection.Name, collection.TotalItems)

		collectionData = append(collectionData, collection)
	}

	return collectionData
}

func extractCategoryAndCollection(filename string) map[string]string {
	pattern := `^(\w+)-(.*?)(\d)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(filename)

	result := make(map[string]string)

	category := matches[1]   // First word
	collection := matches[2] // Everything up to the number
	collection = strings.TrimSuffix(collection, "-")

	if collection == "" {
		collection = category
	}

	result["category"] = category
	result["collection"] = collection

	return result
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
