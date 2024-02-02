package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type PokemonProduct struct {
	url, image, name, price string
}

func writePokemonCSV(pokemonProducts []PokemonProduct) {
	file, err := os.Create("products.csv")

	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	// Write the CSV column headers
	writer.Write([]string{
		"url",
		"image",
		"name",
		"price",
	})

	for _, pokemonProduct := range pokemonProducts {
		pokemonRecord := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		writer.Write(pokemonRecord)
	}

	defer writer.Flush()
}

func scrapePokemonProducts() []PokemonProduct {
	c := colly.NewCollector()
	var pokemonProducts []PokemonProduct

	c.OnHTML("li.product", func(pokemonListElement *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		// Scrape the data we are interested of
		pokemonProduct.url = pokemonListElement.ChildAttr("a", "href")
		pokemonProduct.image = pokemonListElement.ChildAttr("img", "src")
		pokemonProduct.name = pokemonListElement.ChildText("h2")
		pokemonProduct.price = pokemonListElement.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	// Start scraping
	c.Visit("https://scrapeme.live/shop/")

	return pokemonProducts
}

func main() {
	// Scrape pokemon data and store findings to a struct
	pokemonProducts := scrapePokemonProducts()

	// Write CSV with the gathered data
	writePokemonCSV(pokemonProducts)
}
