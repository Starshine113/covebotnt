package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/starshine-sys/covebotnt/structs"
)

func getConfig() (config structs.BotConfig, err error) {
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		sampleConf, err := ioutil.ReadFile("config.sample.toml")
		if err != nil {
			return config, err
		}
		err = ioutil.WriteFile("config.toml", sampleConf, 0644)
		if err != nil {
			return config, err
		}
		fmt.Println("config.toml was not found, created sample configuration.")
		os.Exit(1)
		return config, nil
	}
	configFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(configFile, &config)
	fmt.Println("Loaded configuration file.")
	return config, err
}
