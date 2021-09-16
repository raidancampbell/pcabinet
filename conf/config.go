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
	TestCPU bool
}
type Service struct {
	Name     string
	Endpoint string
}

var C Conf

func Initialize() {
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFilename)
	viper.AddConfigPath("$HOME/.config/pcabinet/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Failed to read in configuration from %s, '%s'", configFilename, err.Error())
	}
	if strings.HasPrefix(C.OutputBasedir, "$HOME") {
		homedir, _ := os.UserHomeDir()
		C.OutputBasedir = strings.Replace(C.OutputBasedir, "$HOME", homedir, 1)
	}

	// tolerable failure, we can fail more thoroughly later with a better message
	os.Mkdir(C.OutputBasedir, 0755)
}
