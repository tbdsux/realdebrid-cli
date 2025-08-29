package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetApiKey() (string, error) {
	apiKey := viper.GetString("apiKey")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set. Please run 'rd config init' to set it.")
	}

	return apiKey, nil
}
