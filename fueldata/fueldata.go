package fueldata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/poolski/fueltracker/types"
)

var endpoints = []string{
	"https://applegreenstores.com/fuel-prices/data.json",
	"https://fuelprices.asconagroup.co.uk/newfuel.json",
	"https://storelocator.asda.com/fuel_prices_data.json",
	"https://www.bp.com/en_gb/united-kingdom/home/fuelprices/fuel_prices_data.json",
	"https://fuelprices.esso.co.uk/latestdata.json",
	"https://jetlocal.co.uk/fuel_prices_data.json",
	"https://www.morrisons.com/fuel-prices/fuel.json",
	"https://moto-way.com/fuel-price/fuel_prices.json",
	"https://fuel.motorfuelgroup.com/fuel_prices_data.json",
	"https://www.rontec-servicestations.co.uk/fuel-prices/data/fuel_prices_data.json",
	"https://api.sainsburys.co.uk/v1/exports/latest/fuel_prices_data.json",
	"https://www.sgnretail.uk/files/data/SGN_daily_fuel_prices.json",
	"https://www.shell.co.uk/fuel-prices-data.html",
}

func GetAllFuelPrices() ([]types.Station, error) {
	res := []types.Station{}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	for _, ep := range endpoints {
		epResponse := &types.FuelDataResponse{}
		fetchWithCaching(client, ep, epResponse)

		res = append(res, epResponse.Stations...)

	}
	return res, nil
}

func fetchWithCaching(client http.Client, ep string, epResponse *types.FuelDataResponse) {
	// Generate a hash of the endpoint to use as a cache key
	hash := sha256.Sum256([]byte(ep))
	hashString := hex.EncodeToString(hash[:])
	cachePath := fmt.Sprintf("cache/%s.json", hashString)

	// Check if the cache exists and is valid. Return the cached data if it is.
	if cacheData, err := os.ReadFile(cachePath); err == nil {
		var cacheItem CacheItem
		if err := json.Unmarshal(cacheData, &cacheItem); err == nil {
			if cacheItem.IsValid() {
				fmt.Printf("using cached data for %s - %v remaining\n", ep, time.Until(cacheItem.Timestamp.Add(time.Hour)))
				err = json.Unmarshal(cacheItem.Data, &epResponse)
				if err != nil {
					fmt.Printf("error unmarshaling response from %s: %v\n", ep, err)
				}
				return
			}
		}
	}

	// If the cache doesn't exist or is invalid, fetch the data and cache it
	readEndpointWithCache(client, ep, cachePath, epResponse)
}

func readEndpointWithCache(client http.Client, ep string, cachePath string, epResponse *types.FuelDataResponse) {
	// Create the cache directory if it doesn't exist
	if _, err := os.Stat("cache"); os.IsNotExist(err) {
		if err = os.Mkdir("cache", 0755); err != nil {
			fmt.Printf("error creating cache directory: %v\n", err)
		}
	}

	// Fetch the data
	fmt.Printf("fetching data from %s\n", ep)
	resp, err := client.Get(ep)
	if err != nil {
		fmt.Printf("error calling %s: %v\n", ep, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response from %s: %v\n", ep, err)
	}

	// Cache the response
	cacheItem := CacheItem{
		Timestamp: time.Now(),
		Data:      body,
	}

	if data, err := json.Marshal(cacheItem); err == nil {
		if err := os.WriteFile(cachePath, data, 0644); err != nil {
			fmt.Printf("error writing cache file: %v\n", err)
		}
	}

	err = json.Unmarshal(cacheItem.Data, &epResponse)
	if err != nil {
		fmt.Printf("error unmarshaling cached response from %s: %v\n", ep, err)
	}
}

func PrintFuelPrices(records []*types.SpecificFuelPrice) {
	// Set up the table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Site ID", "Brand", "Postcode", "Fuel Type", "Price", "Last Recorded At"})
	var tblData [][]string

	// Populate the table
	for _, r := range records {
		tblData = append(tblData, []string{r.SiteID, r.Brand, r.Postcode, r.FuelTypeCode, fmt.Sprintf("%f", r.Price), r.RecordedAt})
	}

	// Draw the table
	for _, v := range tblData {
		table.Append(v)
	}
	table.Render()
}

type CacheItem struct {
	Timestamp time.Time
	Data      []byte
}

func (c *CacheItem) IsValid() bool {
	if time.Since(c.Timestamp) < time.Hour {
		return true
	}
	return false
}
