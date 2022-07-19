package fueldata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
	BaseURL    string
	APIKey     string
	titleCaser cases.Caser
}

func New(APIKey string) *FuelData {
	return &FuelData{
		BaseURL:    "https://uk1.ukvehicledata.co.uk",
		APIKey:     APIKey,
		titleCaser: cases.Title(language.English, cases.NoLower),
	}
}

type QueryOpts struct {
	Postcode string
	FuelType string
	Location string
}

func (c *FuelData) doAPICall(opts QueryOpts) (*types.FuelDataResponse, error) {
	u, _ := url.Parse(c.BaseURL)

	u.Path = fuelPriceEndpoint

	q := u.Query()
	q.Set("v", "2")
	q.Set("api_nullitems", "1")
	q.Set("auth_apikey", c.APIKey)
	q.Set("key_POSTCODE", strings.ToUpper(opts.Postcode))

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

// GetFuelPrices takes a Postcode and a FuelType to show the stations
// which sell that fuel in the search radius for Postcode.
func (c *FuelData) GetFuelPrices(opts QueryOpts) ([]*types.SpecificFuelPrice, error) {
	var prices []*types.SpecificFuelPrice
	if opts.FuelType == "" {
		return nil, errors.New("please specify fuel type")
	}

	// Title case the fuel type for matching on later
	opts.FuelType = c.titleCaser.String(opts.FuelType)

	fd, err := c.doAPICall(opts)
	if err != nil {
		return nil, err
	}

	for _, stn := range fd.DataItems.FuelStationDetails.FuelStationList {
		// If the Location query param is set, skip through the list until we
		// find a fuel station that matches.
		if opts.Location != "" {
			if stn.Name != opts.Location {
				continue
			}
		}

		switch opts.FuelType {
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
	if len(prices) == 0 {
		prices = append(prices, &types.SpecificFuelPrice{
			Station:    "NOTHING FOUND",
			FuelType:   "NOTHING FOUND",
			Price:      0,
			RecordedAt: "",
		})
	}
	return prices, nil
}

func filterPriceByFuel(stn types.FuelStation, ft string) *types.SpecificFuelPrice {
	sfp := &types.SpecificFuelPrice{}
	timeFormat := "1/2/2006 3:04:05 PM"
	for _, fp := range stn.FuelPriceList {
		timestamp, err := time.Parse(timeFormat, fp.LatestRecordedPrice.TimeRecorded)
		if err != nil {
			log.Println(err)
		}
		if fp.FuelType == ft {
			sfp.Station = stn.Name
			sfp.FuelType = ft
			sfp.Price = fp.LatestRecordedPrice.InGbp
			sfp.RecordedAt = timestamp.Local().Format("02/01/2006")
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
