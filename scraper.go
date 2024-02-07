package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type CondoType struct {
	address         string // Viputie 11 A 1
	squareFootage   string // 93,5 m2
	sizeDescription string // 4H+K+S
	buildingType    string // Paritalo
	url             string // https://www.asuntosaatio.fi/asumisoikeusasunnot/espoo/lippajarvi/viputie-11/asunto-a-1/
}

func writeCSV(condos []CondoType) {
	file, err := os.Create("products.csv")

	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	// Write the CSV column headers
	writer.Write([]string{
		"address",
		"squareFootage",
		"sizeDescription",
		"buildingType",
		"url",
	})

	for _, condo := range condos {
		condoRecord := []string{
			condo.address,
			condo.squareFootage,
			condo.sizeDescription,
			condo.buildingType,
			condo.url,
		}

		writer.Write(condoRecord)
	}

	defer writer.Flush()
}

func scrapeCondos(url string) []CondoType {
	var availableCondos []CondoType
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	c.OnHTML("main.py-9.pb-14", func(condoElement *colly.HTMLElement) {
		condo := CondoType{}

		// Scrape the data we are interested of
		condo.address = condoElement.ChildText("h4.text-3xl.mb-4")
		dirtySquareFootageValue := condoElement.ChildText("span.font-bold.whitespace-nowrap")
		// Clean up additional square meter suffix
		squareFootageWithoutSuffix := strings.TrimSuffix(dirtySquareFootageValue, " m2")
		condo.squareFootage = strings.Replace(squareFootageWithoutSuffix, ",", ".", 1)

		// "sizeDescription" and "buildingType" are on the same line as one string
		description := condoElement.ChildText("span.font-normal.text-right")
		descriptionStrings := strings.Split(description, ",")
		condo.sizeDescription = strings.TrimSpace(descriptionStrings[0])
		condo.buildingType = strings.TrimSpace(descriptionStrings[1])
		// TODO: Why no URL?
		condo.url = condoElement.ChildAttr("a.btn-green-400", "href")

		availableCondos = append(availableCondos, condo)
	})

	// Start scraping
	c.Visit(url)

	return availableCondos
}

func main() {
	// Scrape data and store findings to a struct
	availableCondos := scrapeCondos("https://www.asuntosaatio.fi/asunnot/etsi-asuntoa/?cities=Espoo&minSquareMeters=90&buildingTypes=Paritalo,Erillistalo,Rivitalo&roomTypes=4,5,6-99&type=AsoFilters")

	// Write CSV with the gathered data
	writeCSV(availableCondos)
}
