package nominatim

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/poolski/fueltracker/types"
)

const reverseSearchURL = "https://nominatim.openstreetmap.org/reverse"
const forwardSearchURL = "https://nominatim.openstreetmap.org/search"

// PostcodeToCoords takes a UK postcode and returns the latitude and longitude
// of the location. The OSM API responds with a string for both latitude and longitude
// but it is converted to a float64 for ease of use.
func PostcodeToCoords(postcode string) (float64, float64, error) {
	url, err := url.Parse(forwardSearchURL)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing URL: %v", err)
	}
	q := url.Query()
	q.Set("format", "json")
	q.Set("postalcode", postcode)
	q.Set("country", "United Kingdom") // This is a UK only app
	q.Set("limit", "1")
	url.RawQuery = q.Encode()

	res, err := http.Get(url.String())
	if err != nil {
		return 0, 0, fmt.Errorf("error getting response from %s: %v", url.String(), err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading response body: %v", err)
	}

	// The OSM API returns a response object that's nested inside an array, but we only want the first
	// element of the array
	var resp []types.OSMForwardSearchResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, 0, fmt.Errorf("error unmarshaling response body: %v", err)
	}
	lat, err := strconv.ParseFloat(resp[0].Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing latitude: %v", err)
	}

	lon, err := strconv.ParseFloat(resp[0].Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing longitude: %v", err)
	}

	return lat, lon, nil
}
