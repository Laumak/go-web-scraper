package geolocator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
	"web-scraper/scraper"
)

type BatchGeocodeType []struct {
	Query struct {
		Text string `json:"text"`
	} `json:"query"`
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func makeRequest(method string, url string, data []byte) *http.Response {
	fmt.Println(method, data)
	client := &http.Client{}

	// Get content from batch request
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	respDump, _ := httputil.DumpResponse(response, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RESPONSE:\n%s", string(respDump))

	return response
}

func getBatchResultsUrl() string {
	// https://myprojects.geoapify.com/api/YzuOvi4LdmFMK7jCYQFL/keys
	apiKey := os.Getenv("GEOAPIFY_API_KEY")
	url := "https://api.geoapify.com/v1/batch/geocode/search?apiKey=" + apiKey

	// https://gobyexample.com/json
	slice := []string{"Nygrannaksentie 3 G 12", "Vanha Sveinsintie 6 E 15"}
	requestBody, _ := json.Marshal(slice)

	response := makeRequest("POST", url, []byte(requestBody))
	responseBody, _ := io.ReadAll(response.Body)

	type Response struct {
		Id     string `json:"id"`
		Status string `json:"status"`
		Url    string `json:"url"`
	}
	var r Response
	json.Unmarshal(responseBody, &r)

	fmt.Println("Got batch process URL: ", r.Url)

	return r.Url
}

func getBatchResultsFromUrl(url string) any {
	// Get content from batch request
	response := makeRequest("GET", url, nil)
	responseBody, _ := io.ReadAll(response.Body)

	var geocodeResults BatchGeocodeType
	jsonMarshalError := json.Unmarshal(responseBody, &geocodeResults)
	if jsonMarshalError != nil {
		fmt.Println("Unmarshal error: ", jsonMarshalError)
	}

	fmt.Println("RESULTS: ", geocodeResults)

	return geocodeResults
}

func readCondosJsonFileContents(filename string) []scraper.CondoType {
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var content []scraper.CondoType
	err = json.Unmarshal(contents, &content)
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func addGeocodeInfoToCondoJson(condos []scraper.CondoType, geocodeResults any) {
	fmt.Println("HERE-------------------", geocodeResults)
	for i := range condos {
		// TODO: Use dynamic values from
		if condos[i].Address == "Nygrannaksentie 3 G 12" {
			condos[i].Lat = "11"
			condos[i].Lon = "11"
		}
	}

	scraper.WriteCondosJSON(condos)
}

func Geocode() {
	batchResultsUrl := getBatchResultsUrl()

	// Wait for batch results to be ready.
	// We could also wait for "pending" "status" to go away in the batch results
	// -> results are available. This could be polled in certain intervals (1 second?).
	time.Sleep(5 * time.Second)

	geocodeResults := getBatchResultsFromUrl(batchResultsUrl)

	oldCondoJson := readCondosJsonFileContents("condos.json")

	addGeocodeInfoToCondoJson(oldCondoJson, geocodeResults)
}
