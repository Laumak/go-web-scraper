# Golang web scraper & ASO locator frontend

Webscraper library used: Gocolly - https://go-colly.org/

Geocoding API: Geoapify - https://myprojects.geoapify.com/api/YzuOvi4LdmFMK7jCYQFL/keys

Google maps API key: https://console.cloud.google.com/google/maps-apis/credentials?project=aso-locator

## Instructions

For running the go application (web scraper & geolocator), `GEOAPIFY_API_KEY` has to be set in the environment for the application to work. See above for the key location. Command to run the app:

```go
GEOAPIFY_API_KEY=KEY_HERE go run main.go
```

For running the frontend, a (git ignored) `env.js` file must be added to the "frontend" directory with the following contents:

```js
const MAPS_API_KEY = "MAPS_API_KEY_HERE";
const MAP_ID = "MAP_ID_HERE";
```

## TODO

1. Create server for JSON content to be served for the frontend (so that the condo data doesn't have to be copied to FE code)
