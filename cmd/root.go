/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ~/.config/fueltracker/config.json)")
	rootCmd.PersistentFlags().StringP("postcode", "p", "", "postcode to look up fuel prices for e.g. 'AB123XY'")
	if err := rootCmd.MarkPersistentFlagRequired("postcode"); err != nil {
		log.Fatal(err)
	}
	rootCmd.PersistentFlags().StringP("fuel", "f", "unleaded", "(optional) specific fuel type to show prices for")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not read user's config directory: %v", err)
	}
	sep := string(filepath.Separator)
	confDir := homeDir + sep + ".config"
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(confDir + "/fueltracker/")
		viper.SetConfigName("config")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}
