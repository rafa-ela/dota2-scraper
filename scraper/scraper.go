package scraper

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/dota2-scraper/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	packageName = "scraper"
)

// ScrapePlayerIds gets the player IDs, and href links of all the Dota players,
// and writes them to a JSON file.
func ScrapePlayerIds(log logrus.FieldLogger) bool {
	playerMap := make(map[string]api.PlayerLink, 0)

	log = log.WithFields(logrus.Fields{
		"package_name": packageName,
	})

	//TODO: read from the .env file
	allPLayersURL := "https://liquipedia.net/dota2/Players_(all)"

	// Request the HTML page.
	res, err := http.Get(allPLayersURL)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "requesting to url %s", allPLayersURL))
		return false
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return false
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(errors.Wrap(err, "loading the HTML document"))
		return false
	}

	// Find the href, and player IDs
	doc.Find(".sortable.wikitable.smwtable tbody tr").Each(func(i int, s *goquery.Selection) {
		tdNode := s.Find("td").Children()
		player := api.PlayerLink{}

		if tdNode.Size() > 0 {
			idNode := tdNode.Get(1)
			player.Href = idNode.Attr[0].Val
			player.ID = idNode.FirstChild.Data

			playerMap[player.ID] = player
		}
	})

	n, err := WriteDataToFileAsJSON(playerMap, "ids.json")
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Infof("ScrapePlayerIds printed %d bytes to %s.json\n", n, "ids")
	log.Infof("ScrapePlayerIds wrote json %d records to ids.json\n", len(playerMap))

	return true
}

// ScrapePlayerInfo reads the player info
func ScrapePlayerInfo() bool {

	return true
}
