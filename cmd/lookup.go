/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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
	RunE: func(cmd *cobra.Command, args []string) error {
		c := fueldata.New(viper.GetString("ukvd_api_key"))

		postcode, _ := cmd.Flags().GetString("postcode")
		fuel, _ := cmd.Flags().GetString("fuel")
		station, _ := cmd.Flags().GetString("station")

		opts := fueldata.QueryOpts{
			FuelType: fuel,
		}

		fuelData, err := c.FetchFuelDataForPostcode(postcode)
		if err != nil {
			return err
		}
		records := c.ShowSpecificFuelPrices(fuelData, opts)
		if station != "" {
			records = c.GetPriceForLocation(records, station)
		}
		fueldata.PrintFuelPrices(records)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
