/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/poolski/fueltracker/config"
	"github.com/poolski/fueltracker/fueldata"
	"github.com/poolski/fueltracker/sheets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write results to a Google Sheets spreadsheet",
	Long:  `Writes the fuel prices for a specific fuel station out to a Google Sheets spreadsheet`,
	RunE: func(cmd *cobra.Command, args []string) error {
		postcode, _ := cmd.Flags().GetString("postcode")
		fuel, _ := cmd.Flags().GetString("fuel")
		station, _ := cmd.Flags().GetString("station")

		log.Printf("fetching %s fuel prices for %s...", fuel, station)

		cfg := &config.GoogleConfig{
			CredentialsPath: viper.GetString("google.service_account"),
			SpreadsheetID:   viper.GetString("google.spreadsheet_id"),
			WorksheetRange:  viper.GetString("google.worksheet_range"),
		}

		opts := fueldata.QueryOpts{
			FuelType: fuel,
		}

		sheets, err := sheets.New(cfg)
		if err != nil {
			return err
		}

		c := fueldata.New(viper.GetString("ukvd_api_key"))

		fuelData, err := c.FetchFuelDataForPostcode(postcode)
		if err != nil {
			return err
		}

		records := c.ShowSpecificFuelPrices(fuelData, opts)
		if station != "" {
			records = c.GetPriceForLocation(records, station)
		}

		if err := sheets.WriteRecordToSpreadsheet(records[0]); err == nil {
			log.Println("successfully written latest price to Google Sheets")
			return nil
		} else {
			return err
		}
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	if err := rootCmd.MarkPersistentFlagRequired("station"); err != nil {
		log.Println(err)
	}
}
