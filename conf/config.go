package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/andrejsstepanovs/go-litellm/conf/connections"
)

type App struct {
	Connections connections.Config
}

func (a *App) Validate() error {
	var errs []error

	err := a.Connections.Validate()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	var finalErr error
	for _, e := range errs {
		finalErr = fmt.Errorf("%w err: %w", finalErr, e)
	}

	if finalErr != nil {
		return fmt.Errorf("app validation failed: %w", finalErr)
	}
	return nil
}

func Load() (App, error) {
	//Step 1: Set the config file name and type
	viper.SetConfigName("bobik") // Name of the config file (without extension)
	viper.SetConfigType("env")   // Type of the config file

	// Step 2: Add search paths for the config file
	// First, look in the current directory
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../..")

	// Fallback to the user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return App{}, fmt.Errorf("error getting user home directory: %w", err)
	}
	viper.AddConfigPath(home)

	// Step 3: Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found in current directory or home directory")
		} else {
			log.Println("Error reading config file:", err)
		}
		return App{}, fmt.Errorf("error reading config file: %w", err)
	}

	conn, err := connections.New()
	if err != nil {
		return App{}, fmt.Errorf("error creating connections: %w", err)
	}
	app := App{
		Connections: conn,
	}

	err = app.Validate()
	if err != nil {
		return App{}, fmt.Errorf("error validating config: %w", err)
	}

	return app, nil
}
