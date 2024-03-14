package sheets

import (
	"context"
	"fmt"

	"github.com/poolski/fueltracker/config"
	"github.com/poolski/fueltracker/types"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GSheets struct {
	Config  config.GoogleConfig
	Service *sheets.Service
}

func New(cfg *config.GoogleConfig) (*GSheets, error) {
	ctx := context.Background()
	client := &GSheets{
		Config: *cfg,
	}
	svc, err := sheets.NewService(ctx, option.WithCredentialsFile(cfg.CredentialsPath), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return nil, fmt.Errorf("unable to create spreadsheets client: %v", err)
	}

	client.Service = svc
	return client, nil
}

func (s *GSheets) Write(rec *types.SpecificFuelPrice) error {
	spreadsheetID := s.Config.SpreadsheetID
	writeRange := s.Config.WorksheetRange

	var vr sheets.ValueRange

	vals := []interface{}{rec.RecordedAt, rec.FuelTypeCode, rec.Price}
	vr.Values = append(vr.Values, vals)

	_, err := s.Service.Spreadsheets.Values.Append(
		spreadsheetID, writeRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("failed to write to spreadsheet: %v", err)
	}
	return nil
}
