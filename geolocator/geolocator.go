package geolocator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Geocode() {
	// https://myprojects.geoapify.com/api/YzuOvi4LdmFMK7jCYQFL/keys
	apiKey := os.Getenv("GEOAPIFY_API_KEY")
	url := "https://api.geoapify.com/v1/batch/geocode/search?apiKey=" + apiKey

	client := &http.Client{}
	// https://gobyexample.com/json
	slice := []string{"Viputie 11", "Juvanpuistonkuja 2"}
	requestBody, _ := json.Marshal(slice)

	request, requestError := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	if requestError != nil {
		fmt.Println(requestError)
		return
	}

	response, responseError := client.Do(request)
	if responseError != nil {
		fmt.Println(responseError)
		return
	}
	defer response.Body.Close()
	responseBody, _ := io.ReadAll(response.Body)

	type Response struct {
		Id     string `json:"id"`
		Status string `json:"status"`
		Url    string `json:"url"`
	}

	var r Response

	json.Unmarshal(responseBody, &r)

	fmt.Println("Got batch process URL: ", r.Url)

	time.Sleep(5 * time.Second)

	// Get content from batch request
	request2, requestError2 := http.NewRequest("GET", url, bytes.NewBuffer([]byte(r.Url)))
	if requestError2 != nil {
		fmt.Println(requestError2)
		return
	}

	response2, responseError2 := client.Do(request2)
	if responseError2 != nil {
		fmt.Println(responseError2)
		return
	}
	defer response2.Body.Close()
	responseBody2, _ := io.ReadAll(response2.Body)

	type BatchGeocodeType []struct {
		Query struct {
			Text string `json:"text"`
		} `json:"query"`
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}

	var unMarshaledBatchResultResponse BatchGeocodeType

	jsonMarshalError := json.Unmarshal(responseBody2, &unMarshaledBatchResultResponse)
	if jsonMarshalError != nil {
		fmt.Println(jsonMarshalError)
		return
	}

	fmt.Println(unMarshaledBatchResultResponse)
}
