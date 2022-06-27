package config

type GoogleConfig struct {
	CredentialsPath string `mapstructure:"credentials_path"`
	SpreadsheetID   string `mapstructure:"spreadsheet_id"`
	WorksheetRange  string `mapstructure:"worksheet_range"`
}
