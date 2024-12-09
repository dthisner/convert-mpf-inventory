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

func GetCollections() []models.CollectionData {
	smallsCollectionID := map[int]string{
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

	lightningCollectionsID := map[int]string{
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

	_ = lightningCollectionsID

	log.Print("Getting ALL the collections")

	var collectionData []models.CollectionData
	for id, name := range smallsCollectionID {
		var collection models.CollectionData
		start := 0
		end := 200

		MPF, err := getCollectionData(id, start, end)
		if err != nil {
			log.Print(err)
		}

		collection.MRP_DATA = MPF
		collection.ID = id
		collection.Category = "lightning"
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
