package main

import (
	"github.com/dota2-scraper/scraper"
	"github.com/sirupsen/logrus"
)

const (
	packageName = "main"
	serviceName = "dota2scraper"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"service": serviceName,
		"package": packageName,
	})

	succesful := scraper.ScrapePlayerIds(logger)

	if !succesful {
		logger.Fatal("unable to scrape player ids")
	}
}
