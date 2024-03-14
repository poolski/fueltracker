package types

import (
	"encoding/json"
	"strconv"
)

// OSMForwardSearchResponse is the response from the OSM API when doing a forward search
type OSMForwardSearchResponse struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	AddressType string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	BoundingBox []string `json:"boundingbox"`
}

// Fuel Station Responses

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Station struct {
	SiteID   string             `json:"site_id"`
	Brand    string             `json:"brand"`
	Address  string             `json:"address"`
	Postcode string             `json:"postcode"`
	Location Location           `json:"location"`
	Prices   map[string]float64 `json:"prices"`
}

type FuelDataResponse struct {
	LastUpdated string    `json:"last_updated"`
	Stations    []Station `json:"stations"`
}

type SpecificFuelPrice struct {
	SiteID       string
	Brand        string
	Postcode     string
	FuelTypeCode string
	Price        float64
	RecordedAt   string
	MonthYear    string
}

// UnmarshalJSON for Location to handle both float64 and string types for latitude and longitude.
// This is necessary become some clown didn't read the official spec here: https://assets.publishing.service.gov.uk/media/64d4ac7b5cac65000dc2dd1a/A._Appendix_A.pdf and decided to use strings for numbers.
// Looking at you, Morrisons.
func (l *Location) UnmarshalJSON(data []byte) error {
	// Intermediate representation of Location, using interface{} for latitude and longitude.
	aux := struct {
		Latitude  interface{} `json:"latitude"`
		Longitude interface{} `json:"longitude"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Attempt to handle latitude as float64 or string.
	if lat, ok := aux.Latitude.(float64); ok {
		l.Latitude = lat
	} else if latStr, ok := aux.Latitude.(string); ok {
		var err error
		l.Latitude, err = strconv.ParseFloat(latStr, 64)
		if err != nil {
			return err
		}
	}

	// Attempt to handle longitude as float64 or string.
	if lon, ok := aux.Longitude.(float64); ok {
		l.Longitude = lon
	} else if lonStr, ok := aux.Longitude.(string); ok {
		var err error
		l.Longitude, err = strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return err
		}
	}

	return nil
}
