package conf

import (
	"log"
	"os"
	"strings"
)
import "github.com/spf13/viper"

const configFilename = "pcabinet_config.yml"

type Conf struct {
	Services []Service
	OutputBasedir string `mapstructure:"output_basedir"`
}
type Service struct {
	Name     string
	Endpoint string
}

var OutputBasedir string

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
	if strings.HasPrefix(c.OutputBasedir, "$HOME") {
		homedir, _ := os.UserHomeDir()
		c.OutputBasedir = strings.Replace(c.OutputBasedir, "$HOME", homedir, 1)
	}
	OutputBasedir = c.OutputBasedir
	os.UserHomeDir()

	// tolerable failure, we can fail more thoroughly later with a better message
	os.Mkdir(OutputBasedir, 0755)

	return &c
}
