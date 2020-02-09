package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/dota2-scraper/api"
)

// ScrapePlayerIds gets the player IDs, and href links of all the Dota players,
// and writes them to a JSON file.
func ScrapePlayerIds() {
	// Request the HTML page.
	res, err := http.Get("https://liquipedia.net/dota2/Players_(all)")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	playerMap := make(map[string]api.PlayerLink, 0)

	// Find the items
	doc.Find(".sortable.wikitable.smwtable tbody tr").Each(func(i int, s *goquery.Selection) {
		tdNode := s.Find("td").Children()
		player := api.PlayerLink{}

		if tdNode.Size() > 0 {
			idNode := tdNode.Get(1)
			player.Href = idNode.Attr[0].Val
			player.ID = idNode.FirstChild.Data

			playerMap[player.ID] = player

			//marshal to JSON and write to file
		}
	})

	n, err := WriteDataToFileAsJSON(playerMap, "ids.json")
	if err != nil {
		log.Fatal(err)
	}

	//TODO: use logger instead
	fmt.Printf("ScrapePlayerIds printed %d bytes to %s.json\n", n, "ids")
	fmt.Printf("ScrapePlayerIds wrote json %d records to ids.json\n", len(playerMap))
}

// WriteDataToFileAsJSON is taken from: https://stackoverflow.com/a/57192522
func WriteDataToFileAsJSON(data interface{}, filedir string) (int, error) {
	//write data as buffer to json encoder
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return 0, err
	}
	file, err := os.OpenFile(filedir, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return 0, err
	}
	n, err := file.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}
