package initial

import (
	"fmt"
	"lesson/simple-bank/config"

	"github.com/spf13/viper"
)

func LoadingConfig(path string) (config config.Config, err error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)               // optionally look for config in the working directory
	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil { // Handle errors unmarshaling the config
	    panic(fmt.Errorf("fatal error config file: %w", err))
	}
    return
}