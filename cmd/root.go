package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fueltracker",
	Short: "A CLI tool to look up fuel prices locally",
	Long:  `This tool looks up fuel prices from the UK Vehicle Data API`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check for config file
		if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
			fmt.Printf("config file not found: %v\n", err)
			configCmd.Root()
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cobra.OnInitialize(initConfig)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not read user's config directory: %v", err)
	}
	sep := string(filepath.Separator)
	confDir := homeDir + sep + ".config"
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", confDir+"/fueltracker/config.json", "config file")

	rootCmd.PersistentFlags().StringP("postcode", "p", "", "postcode to look up fuel prices for e.g. 'AB123XY'")
	rootCmd.PersistentFlags().StringP("fuel", "f", "E10", "(optional) fuel code to look up. Valid options are 'E10' (Unleaded), 'E5' (Premium unleaded), 'B7' (Diesel), 'SDV' (Premium diesel)")
	rootCmd.PersistentFlags().StringP("site-id", "s", "", "(optional) site ID to look up. If this is set, the postcode flag is ignored")
	rootCmd.PersistentFlags().StringP("radius", "r", "5.0", "(optional) radius in km to search around the postcode")

	// Check if site ID is set, ignore postcode flag
	if siteID := rootCmd.PersistentFlags().Lookup("site-id").Value.String(); siteID != "" {
		if err := rootCmd.PersistentFlags().SetAnnotation("postcode", cobra.BashCompOneRequiredFlag, []string{"false"}); err != nil {
			log.Fatal(err)
		}
	} else {
		// Require postcode flag
		if err := rootCmd.PersistentFlags().SetAnnotation("postcode", cobra.BashCompOneRequiredFlag, []string{"true"}); err != nil {
			log.Fatal(err)
		}
	}

	// The `generate` command doesn't need the postcode flag
	if err := configCmd.InheritedFlags().SetAnnotation("postcode", cobra.BashCompOneRequiredFlag, []string{"false"}); err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		configCmd.Root()
	}
}
