package mpf

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"export-mountpf-inventory/models"
)

var (
	API_KEY string = "311bdd33-497e-4579-8468-75257954f7b8"
)

// Loopa igenom Collections för att få ALLA collections, MPF
// Spara ut dem som CSV

func getCollectionData(id, start, end int) (models.MPF_EXPORT, error) {
	var MPF models.MPF_EXPORT
	log.Printf("Getting for Collection ID: %d from %d to %d", id, start, end)

	url := fmt.Sprintf("https://svc-1001-usf.hotyon.com/search?q=&apiKey=%s&locale=en&collection=%d&skip=%d&take=%d&sort=title", API_KEY, id, start, end)
	// log.Printf("The Log: \"%s\"", url)

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
