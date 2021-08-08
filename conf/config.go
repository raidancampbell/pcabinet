package conf

import (
	"log"
)
import "github.com/spf13/viper"
const configFilename = "config.yml"

type Conf struct {
	Services map[string]Service
}
type Service struct {
	Endpoint string
}

func Initialize() *Conf {
	var c Conf
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFilename)
	err := viper.ReadInConfig()
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("Failed to read in configuration from %s, '%s'", configFilename, err.Error())
	}

	return &c
}
