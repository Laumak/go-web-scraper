package scraper

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type CondoType struct {
	Address         string `json:"address"`         // Viputie 11 A 1
	SquareFootage   string `json:"squareFootage"`   // 93,5 m2
	SizeDescription string `json:"sizeDescription"` // 4H+K+S
	BuildingType    string `json:"buildingType"`    // Paritalo
	Url             string `json:"url"`             // https://www.asuntosaatio.fi/asumisoikeusasunnot/espoo/lippajarvi/viputie-11/asunto-a-1/
}

func writeCSV(condos []CondoType) {
	file, err := os.Create("condos.csv")

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
			condo.Address,
			condo.SquareFootage,
			condo.SizeDescription,
			condo.BuildingType,
			condo.Url,
		}

		writer.Write(condoRecord)
	}

	defer writer.Flush()
}

func writeJSON(availableCondos []CondoType) {
	// Serialize the struct to JSON
	jsonBytes, err := json.MarshalIndent(availableCondos, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the JSON data to a file
	err = os.WriteFile("condos.json", jsonBytes, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func scrapeCondos(url string) []CondoType {
	var availableCondos []CondoType
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	c.OnHTML("main.py-9.pb-14", func(condoElement *colly.HTMLElement) {
		condo := CondoType{}

		// Scrape the data we are interested of
		condo.Address = condoElement.ChildText("h4.text-3xl.mb-4")
		dirtySquareFootageValue := condoElement.ChildText("span.font-bold.whitespace-nowrap")
		// Clean up additional square meter suffix
		squareFootageWithoutSuffix := strings.TrimSuffix(dirtySquareFootageValue, " m2")
		condo.SquareFootage = strings.Replace(squareFootageWithoutSuffix, ",", ".", 1)

		// "sizeDescription" and "BuildingType" are on the same line as one string
		description := condoElement.ChildText("span.font-normal.text-right")
		descriptionStrings := strings.Split(description, ",")
		condo.SizeDescription = strings.TrimSpace(descriptionStrings[0])
		condo.BuildingType = strings.TrimSpace(descriptionStrings[1])
		// TODO: Why no URL?
		condo.Url = condoElement.ChildAttr("a", "href")

		availableCondos = append(availableCondos, condo)
	})

	// Start scraping
	c.Visit(url)

	return availableCondos
}

func Scrape() {
	// Scrape data and store findings to a struct
	availableCondos := scrapeCondos("https://www.asuntosaatio.fi/asunnot/etsi-asuntoa/?cities=Espoo&minSquareMeters=90&buildingTypes=Paritalo,Erillistalo,Rivitalo&roomTypes=4,5,6-99&type=AsoFilters")

	// Write CSV with the gathered data
	writeCSV(availableCondos)

	writeJSON(availableCondos)
}
