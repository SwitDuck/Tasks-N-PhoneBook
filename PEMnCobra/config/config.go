package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadYamlConfig()(*viper.Viper, error){
	localViper := viper.New()
	localViper.SetConfigName("config")
	localViper.SetConfigType("yaml")
	localViper.AddConfigPath(".")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := localViper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok{
			logger.Info("Config file not found, using default setting")
		} else {
			return localViper, fmt.Errorf("Found config file, but encountered an error: %v", err)
		}
	} 
	return localViper, nil
}