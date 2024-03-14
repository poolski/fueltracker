package cmd

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/poolski/fueltracker/fueldata"
	"github.com/poolski/fueltracker/nominatim"
	"github.com/poolski/fueltracker/types"
	"github.com/spf13/cobra"
)

const EarthRadius = 6371.0 // Earth's radius in km
const (
	FuelTypeUnleaded      = "E10"
	FuelTypeSuperUnleaded = "E5"
	FuelTypeDiesel        = "B7"
	FuelTypePremiumDiesel = "SDV"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Look up fuel data for a postcode",
	Long: `lookup --postcode AB123XY
	
Lookup is used to get fuel prices matching specific search criteria:
- postcode
- fuel type
- radius around postcode
Its purpose is to allow you to find a specific fuel price (or prices) for a specific location (site ID)
Once you have a site ID, you can use the write command to write the data to a Google Sheet.`,
	RunE: doLookup,
}

func doLookup(cmd *cobra.Command, args []string) error {
	postcode := cmd.Flags().Lookup("postcode").Value.String()
	radius, err := strconv.ParseFloat(cmd.Flags().Lookup("radius").Value.String(), 64)
	if err != nil {
		return fmt.Errorf("error parsing radius: %v", err)
	}

	fuelCode := cmd.Flags().Lookup("fuel").Value.String()
	switch fuelCode {
	case FuelTypeUnleaded, FuelTypeSuperUnleaded, FuelTypeDiesel, FuelTypePremiumDiesel:
		// Do nothing
	default:
		return fmt.Errorf("invalid fuel code")
	}

	pcLat, pcLon, err := nominatim.PostcodeToCoords(postcode)
	if err != nil {
		return err
	}

	records, err := fueldata.GetAllFuelPrices()
	if err != nil {
		return err
	}

	filteredRecords := []*types.SpecificFuelPrice{}

	for _, r := range records {
		lat, long := r.Location.Latitude, r.Location.Longitude // Asserting the type of lat and long as float64
		if IsInsideCircle(pcLat, pcLon, lat, long, radius) {
			for k, v := range r.Prices {
				// If no fuelCode is specified, then we want to print all fuel prices for the station
				if fuelCode == "" || fuelCode == k {
					filteredRecords = append(filteredRecords, &types.SpecificFuelPrice{
						SiteID:       r.SiteID,
						Brand:        r.Brand,
						Postcode:     r.Postcode,
						FuelTypeCode: k,
						Price:        v,
						RecordedAt:   time.Now().Format("2006-01-02 15:04:05"),
						MonthYear:    time.Now().Format("01/2006"),
					})
				}
			}
		}
	}

	fueldata.PrintFuelPrices(filteredRecords)
	return nil
}

// Haversine calculates the distance between two points on the Earth's surface using the Haversine formula.
// https://en.wikipedia.org/wiki/Haversine_formula
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude and longitude from degrees to radians.
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Calculate differences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Apply Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate distance
	distance := EarthRadius * c

	return distance
}

// IsInsideCircle checks if a point is inside a circle. This is a much simpler way to work out
// if a location is within a certain radius of a postcode.
func IsInsideCircle(centerLat, centerLon, pointLat, pointLon, radiusKm float64) bool {
	distance := Haversine(centerLat, centerLon, pointLat, pointLon)
	return distance <= radiusKm
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
