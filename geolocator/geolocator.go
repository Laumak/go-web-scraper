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
	"strconv"
	"time"
	"web-scraper/scraper"
)

type BatchGeocodeType struct {
	Query struct {
		Text string `json:"text"`
	} `json:"query"`
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type BatchGeocodeTypeResponse []BatchGeocodeType

func makeRequest(method string, url string, data []byte) *http.Response {
	fmt.Println(method, ": ", url)
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

func createSliceOfCondoAddresses(condos []scraper.CondoType) []string {
	var sliceOfStrings []string

	for i := range condos {
		sliceOfStrings = append(sliceOfStrings, condos[i].Address)
	}

	return sliceOfStrings
}

func getBatchResultsUrl(condoJson []scraper.CondoType) string {
	// https://myprojects.geoapify.com/api/YzuOvi4LdmFMK7jCYQFL/keys
	apiKey := os.Getenv("GEOAPIFY_API_KEY")
	url := "https://api.geoapify.com/v1/batch/geocode/search?apiKey=" + apiKey

	// https://gobyexample.com/json
	sliceOfCondoAddresses := createSliceOfCondoAddresses(condoJson)
	requestBody, _ := json.Marshal(sliceOfCondoAddresses)

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

func getBatchResultsFromUrl(url string) BatchGeocodeTypeResponse {
	// Get content from batch request
	response := makeRequest("GET", url, nil)
	responseBody, _ := io.ReadAll(response.Body)

	var geocodeResults BatchGeocodeTypeResponse
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

func addGeocodeInfoToCondoJson(condos []scraper.CondoType, geocodeResults []BatchGeocodeType) {
	for i := range geocodeResults {
		for j := range condos {
			if geocodeResults[i].Query.Text == condos[j].Address {
				condos[j].Lat = strconv.FormatFloat(geocodeResults[i].Lat, 'f', -1, 64)
				condos[j].Lon = strconv.FormatFloat(geocodeResults[i].Lon, 'f', -1, 64)
			}
		}
	}

	scraper.WriteCondosJSON(condos)
}

func Geocode() {
	condoJson := readCondosJsonFileContents("condos.json")
	batchResultsUrl := getBatchResultsUrl(condoJson)

	// Wait for batch results to be ready.
	// We could also wait for "pending" "status" to go away in the batch results
	// -> results are available. This could be polled in certain intervals (1 second?).
	fmt.Println("Waiting for 45 seconds")
	time.Sleep(45 * time.Second)

	geocodeResults := getBatchResultsFromUrl(batchResultsUrl)

	oldCondoJson := readCondosJsonFileContents("condos.json")

	addGeocodeInfoToCondoJson(oldCondoJson, geocodeResults)
}
