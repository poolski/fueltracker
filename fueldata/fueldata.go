package fueldata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/poolski/fueltracker/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const fuelPriceEndpoint = "api/datapackage/FuelPriceData"
const (
	FuelTypeUnleaded      = "Unleaded"
	FuelTypeSuperUnleaded = "Super Unleaded"
	FuelTypeDiesel        = "Diesel"
	FuelTypePremiumDiesel = "Premium Diesel"
)

type FuelData struct {
	BaseURL string
	APIKey  string
	Caser   cases.Caser
}

func New(APIKey string) *FuelData {
	return &FuelData{
		BaseURL: "https://uk1.ukvehicledata.co.uk",
		APIKey:  APIKey,
		Caser:   cases.Title(language.English, cases.NoLower),
	}
}

type QueryOpts struct {
	FuelType string
	Location string
}

func (c *FuelData) FetchFuelDataForPostcode(postcode string) (*types.FuelDataResponse, error) {
	u, _ := url.Parse(c.BaseURL)

	u.Path = fuelPriceEndpoint

	q := u.Query()
	q.Set("v", "2")
	q.Set("api_nullitems", "1")
	q.Set("auth_apikey", c.APIKey)
	q.Set("key_POSTCODE", c.Caser.String(postcode))

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

	if data.Response.StatusCode != "Success" {
		return nil, errors.New(data.Response.StatusMessage)
	}
	return &data.Response, nil
}

// ShowSpecificFuelPrices can be called to get the price of a fuel type
// It returns a list of fuel stations which have the fuel type in a 5km radius of the search post code.
func (c *FuelData) ShowSpecificFuelPrices(fd *types.FuelDataResponse, opts ...QueryOpts) []*types.SpecificFuelPrice {
	var prices []*types.SpecificFuelPrice
	// Parse options if set
	var fuelType string
	for _, opt := range opts {
		if opt.FuelType != "" {
			fuelType = c.Caser.String(opt.FuelType)
		}
	}

	for _, stn := range fd.DataItems.FuelStationDetails.FuelStationList {
		switch fuelType {
		case FuelTypeUnleaded:
			if stn.Features.Fuel.HasUnleaded {
				prices = append(prices, filterPriceByFuel(stn, FuelTypeUnleaded))
			}
		case FuelTypeSuperUnleaded:
			if stn.Features.Fuel.HasSuperUnleaded {
				prices = append(prices, filterPriceByFuel(stn, FuelTypeSuperUnleaded))
			}
		case FuelTypeDiesel:
			if stn.Features.Fuel.HasDiesel {
				prices = append(prices, filterPriceByFuel(stn, FuelTypeDiesel))
			}
		case FuelTypePremiumDiesel:
			if stn.Features.Fuel.HasPremiumDiesel {
				prices = append(prices, filterPriceByFuel(stn, FuelTypePremiumDiesel))
			}
		}
	}
	return prices
}

func (c *FuelData) GetPriceForLocation(prices []*types.SpecificFuelPrice, location string) []*types.SpecificFuelPrice {
	for _, rec := range prices {
		if rec.Station == location {
			return []*types.SpecificFuelPrice{rec}
		}
	}
	return nil
}

func filterPriceByFuel(stn types.FuelStation, ft string) *types.SpecificFuelPrice {
	sfp := &types.SpecificFuelPrice{}
	for _, fp := range stn.FuelPriceList {
		if fp.FuelType == ft {
			sfp.Station = stn.Name
			sfp.FuelType = ft
			sfp.Price = fp.LatestRecordedPrice.InGbp
			sfp.RecordedAt = fp.LatestRecordedPrice.TimeRecorded
		}
	}
	return sfp
}

func PrintFuelPrices(records []*types.SpecificFuelPrice) {
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
