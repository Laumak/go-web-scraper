package main

import (
	"web-scraper/geolocator"
	"web-scraper/scraper"
)

func main() {
	scraper.Scrape()
	geolocator.Geocode()
}
