/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
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
		station, err := cmd.Flags().GetString("station")
		if err != nil {
			return errors.New("station flag required but not set")
		}

		log.Printf("fetching %s fuel prices for %s...", fuel, station)

		cfg := &config.GoogleConfig{
			CredentialsPath: viper.GetString("google.credentials_path"),
			SpreadsheetID:   viper.GetString("google.spreadsheet_id"),
			WorksheetRange:  viper.GetString("google.worksheet_range"),
		}

		opts := fueldata.QueryOpts{
			Postcode: postcode,
			FuelType: fuel,
			Location: station,
		}

		sheets, err := sheets.New(cfg)
		if err != nil {
			return err
		}

		c := fueldata.New(viper.GetString("ukvd_api_key"))

		records, err := c.GetFuelPrices(opts)
		if err != nil {
			return err
		}

		if err := sheets.Write(records[0]); err == nil {
			log.Println("successfully written latest price to Google Sheets")
			return nil
		} else {
			return err
		}
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.Flags().StringP("station", "s", "", "specific fuel station to show prices for")
	if err := writeCmd.MarkFlagRequired("station"); err != nil {
		log.Fatal(err)
	}
}
