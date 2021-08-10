package conf

import (
	"log"
)
import "github.com/spf13/viper"

const configFilename = "pcabinet_config.yml"

type Conf struct {
	Services []Service
}
type Service struct {
	Name     string
	Endpoint string
}

func Initialize() *Conf {
	var c Conf
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFilename)
	viper.AddConfigPath("$HOME/.config/pcabinet/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("Failed to read in configuration from %s, '%s'", configFilename, err.Error())
	}

	return &c
}
