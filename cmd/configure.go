package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/manifoldco/promptui"
	"github.com/poolski/fueltracker/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "configure",
	Short: "Generate a config file for Viper",
	Long:  `Rather than having to copy the sample config file around, this will help you generate a config file and write it to the default location`,
	RunE:  generateCmdRunE,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func generateCmdRunE(cmd *cobra.Command, args []string) error {
	fmt.Println("Let's generate a config file.")
	config := &config.Config{}
	// Iterate over config struct, extract field names and prompt for values
	for i := 0; i < reflect.TypeOf(*config).NumField(); i++ {
		field := reflect.TypeOf(*config).Field(i)
		// Special case for Google struct
		if field.Name == "Google" {
			for j := 0; j < reflect.TypeOf(config.Google).NumField(); j++ {
				googleField := reflect.TypeOf(config.Google).Field(j)
				value, err := promptForValue(googleField.Name)
				if err != nil {
					return err
				}
				reflect.ValueOf(config).Elem().FieldByName(field.Name).FieldByName(googleField.Name).SetString(value)
			}
			continue
		}
		value, err := promptForValue(field.Name)
		if err != nil {
			return err
		}
		reflect.ValueOf(config).Elem().FieldByName(field.Name).SetString(value)
	}

	fmt.Printf("writing config to %s\n", cfgFile)
	os.Exit(0)

	// Marshal config object to JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON to file
	err = os.WriteFile(cfgFile, configJSON, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func promptForValue(key string) (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("value cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    key,
		Validate: validate,
	}

	return prompt.Run()
}
