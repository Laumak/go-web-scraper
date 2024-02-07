package geolocator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Geocode() {
	// https://myprojects.geoapify.com/api/YzuOvi4LdmFMK7jCYQFL/keys
	apiKey := os.Getenv("GEOAPIFY_API_KEY")
	url := "https://api.geoapify.com/v1/batch/geocode/search?apiKey=" + apiKey
	method := "POST"

	client := &http.Client{}
	// https://gobyexample.com/json
	slice := []string{"Viputie 11", "Juvanpuistonkuja 2"}
	requestBody, _ := json.Marshal(slice)

	request, requestError := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
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

	fmt.Println(r.Url)
}
