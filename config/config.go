package config

type AppConfig struct {
	UKVDApiKey string       `mapstructure:"ukvd_api_key"`
	Google     GoogleConfig `mapstructure:"google"`
}

type GoogleConfig struct {
	CredentialsPath string `mapstructure:"credentials_path"`
	SpreadsheetID   string `mapstructure:"spreadsheet_id"`
	WorksheetRange  string `mapstructure:"worksheet_range"`
}
