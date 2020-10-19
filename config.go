package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/Starshine113/covebotnt/structs"
)

func getConfig() (config structs.BotConfig, err error) {
	configFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(configFile, &config)
	sugar.Infof("Loaded configuration file.")
	return config, err
}
