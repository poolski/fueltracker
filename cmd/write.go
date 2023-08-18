package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/PremiereGlobal/go-deadmanssnitch"
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
	RunE:  doWrite,
}

func doWrite(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("creating google sheets connection: %w", err)
	}

	c := fueldata.New(viper.GetString("ukvd_api_key"))

	records, err := c.GetFuelPrices(opts)
	if err != nil {
		return fmt.Errorf("getting fuel prices: %w", err)
	}

	if err := sheets.Write(records[0]); err != nil {
		return err
	} else {
		// If you don't have a Dead Man's Snitch account, we won't do this.
		if c.SnitchAPIKey != "" {
			dms := deadmanssnitch.NewClient(c.SnitchAPIKey)
			if err := dms.CheckIn(viper.GetString("snitch_id")); err != nil {
				log.Printf("writing to DMS: %v", err)
			}
		}
		log.Println("successfully written latest price to Google Sheets")
		return nil
	}
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.Flags().StringP("station", "s", "", "specific fuel station to show prices for")
	if err := writeCmd.MarkFlagRequired("station"); err != nil {
		log.Fatal(err)
	}
}
