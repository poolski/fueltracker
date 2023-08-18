package cmd

import (
	"github.com/poolski/fueltracker/fueldata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Look up fuel data for a postcode",
	Long:  `lookup AB123XY`,
	RunE:  doLookup,
}

func doLookup(cmd *cobra.Command, args []string) error {
	c := fueldata.New(viper.GetString("ukvd_api_key"))

	postcode, _ := cmd.Flags().GetString("postcode")
	fuel, _ := cmd.Flags().GetString("fuel")
	station, _ := cmd.Flags().GetString("station")

	opts := fueldata.QueryOpts{
		Postcode: postcode,
		FuelType: fuel,
		Location: station,
	}

	records, err := c.GetFuelPrices(opts)
	if err != nil {
		return err
	}

	fueldata.PrintFuelPrices(records)
	return nil
}

func init() {
	rootCmd.AddCommand(lookupCmd)
	lookupCmd.Flags().StringP("station", "s", "", "(optional) specific fuel station to show prices for")
}
