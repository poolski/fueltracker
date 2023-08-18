package config

type GoogleConfig struct {
	CredentialsPath string `mapstructure:"credentials_path"`
	SpreadsheetID   string `mapstructure:"spreadsheet_id"`
	WorksheetRange  string `mapstructure:"worksheet_range"`
}

type Config struct {
	UKVDAPIKey   string       `mapstructure:"ukvd_api_key"`
	SnitchAPIKey string       `mapstructure:"snitch_api_key"`
	SnitchID     string       `mapstructure:"snitch_id"`
	Google       GoogleConfig `mapstructure:"google"`
}
