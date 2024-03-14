package cmd

import (
	"fmt"
	"log"

	"github.com/poolski/fueltracker/config"
	"github.com/poolski/fueltracker/sheets"
	"github.com/poolski/fueltracker/types"
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
	// postcode := cmd.Flags().Lookup("postcode").Value.String()
	fuelCode := cmd.Flags().Lookup("fuel").Value.String()
	station := cmd.Flags().Lookup("site-id").Value.String()

	log.Printf("fetching %s fuel prices for %s...", fuelCode, station)

	cfg := &config.GoogleConfig{
		CredentialsPath: viper.GetString("google.credentials_path"),
		SpreadsheetID:   viper.GetString("google.spreadsheet_id"),
		WorksheetRange:  viper.GetString("google.worksheet_range"),
	}

	sheets, err := sheets.New(cfg)
	if err != nil {
		return fmt.Errorf("creating google sheets connection: %w", err)
	}

	records := []*types.SpecificFuelPrice{{
		Brand:        "foo",
		FuelTypeCode: FuelTypeDiesel,
		Price:        0.0,
		RecordedAt:   "",
	}}
	// if err != nil {
	// 	return fmt.Errorf("getting fuel prices: %w", err)
	// }

	if err := sheets.Write(records[0]); err != nil {
		return err
	} else {
		// If you don't have a Dead Man's Snitch account, we won't do this.
		// if c.SnitchAPIKey != "" {
		// 	dms := deadmanssnitch.NewClient(c.SnitchAPIKey)
		// 	if err := dms.CheckIn(viper.GetString("snitch_id")); err != nil {
		// 		log.Printf("writing to DMS: %v", err)
		// 	}
		// }
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
