package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/poolski/fueltracker/types"
)

const timeFormat = "1/2/2012 3:4:05 PM"

type UKVDClient struct {
	BaseURL string
	APIKey  string
}

func NewClient(APIKey string) *UKVDClient {
	return &UKVDClient{
		BaseURL: "https://uk1.ukvehicledata.co.uk",
		APIKey:  APIKey,
	}
}

func (c *UKVDClient) fetchFuelData(postcode string) (*types.FuelDataResponse, error) {
	u, _ := url.Parse(c.BaseURL)

	u.Path = "api/datapackage/FuelPriceData"

	q := u.Query()
	q.Set("v", "2")
	q.Set("api_nullitems", "1")
	q.Set("auth_apikey", c.APIKey)
	q.Set("key_POSTCODE", postcode)

	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := types.RawAPIResponse{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// showUnleadedFuelPrices can be called to get the price of unleaded fuel
// It returns a list of fuel stations which have the fuel type in a 5km radius of the search post code.
func (c *UKVDClient) showUnleadedFuelPrices(fuelData *types.FuelDataResponse) []types.SpecificFuelPrice {
	var prices []types.SpecificFuelPrice
	for _, stn := range fuelData.DataItems.FuelStationDetails.FuelStationList {
		// We only care about unleaded because fuck diesel.
		// TODO: Add the other fuel types that weird people use.
		if stn.Features.Fuel.HasUnleaded {
			for _, fp := range stn.FuelPriceList {
				if fp.FuelType == "Unleaded" {
					prices = append(prices, types.SpecificFuelPrice{
						Station:    stn.Name,
						FuelType:   "Unleaded",
						Price:      fp.LatestRecordedPrice.InPence,
						RecordedAt: fp.LatestRecordedPrice.TimeRecorded,
					})
				}
			}
		}
	}
	return prices
}

func (c *UKVDClient) printFuelPrices(postcode string) {
	fuelData, err := c.fetchFuelData(postcode)
	if err != nil {
		log.Fatal(err)
	}
	records := c.showUnleadedFuelPrices(fuelData)

	// Set up the table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Location", "Fuel Type", "Price", "Last Recorded At"})
	var tblData [][]string

	// Populate the table
	for _, r := range records {
		tblData = append(tblData, []string{r.Station, r.FuelType, fmt.Sprintf("%f", r.Price), r.RecordedAt})
	}

	// Draw the table
	for _, v := range tblData {
		table.Append(v)
	}
	table.Render()
}

func main() {
	//TODO switch to Cobra's cmd.Execute()
	apiKey, ok := os.LookupEnv("UKVD_API_KEY")
	if !ok {
		log.Fatal("Please set the UKVD_API_KEY environment variable")
	}
	postcode := flag.String("postcode", "ABC123", "The postcode you're looking up fuel prices for")
	flag.Parse()

	if !isFlagPassed("postcode") {
		log.Fatal("Please pass the -postcode flag and a value")
	}

	c := NewClient(apiKey)
	// The API expects capitalized postcodes. Natch.
	c.printFuelPrices(strings.ToUpper(*postcode))
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
