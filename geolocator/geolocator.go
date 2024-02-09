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
)

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
	fmt.Printf("RESPONSE:\n%s", string(respDump))

	return response
}

func getBatchResultsUrl() string {
	// https://myprojects.geoapify.com/api/YzuOvi4LdmFMK7jCYQFL/keys
	apiKey := os.Getenv("GEOAPIFY_API_KEY")
	url := "https://api.geoapify.com/v1/batch/geocode/search?apiKey=" + apiKey

	// https://gobyexample.com/json
	slice := []string{"Viputie 11", "Juvanpuistonkuja 2"}
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

func getBatchResultsFromUrl(url string) {
	// Get content from batch request
	response := makeRequest("GET", url, nil)
	responseBody, _ := io.ReadAll(response.Body)

	type BatchGeocodeType []struct {
		Query struct {
			Text string `json:"text"`
		} `json:"query"`
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
	var geocodeResults BatchGeocodeType
	jsonMarshalError := json.Unmarshal(responseBody, &geocodeResults)
	if jsonMarshalError != nil {
		fmt.Println("Unmarshal error: ", jsonMarshalError)
		return
	}

	fmt.Println("RESULTS: ", geocodeResults)
}

func Geocode() {
	batchResultsUrl := getBatchResultsUrl()

	time.Sleep(5 * time.Second)

	getBatchResultsFromUrl(batchResultsUrl)
}
