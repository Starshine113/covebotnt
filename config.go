package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type botConfig struct {
	Auth struct {
		Token       string `toml:"token"`
		DatabaseURL string `toml:"database_url"`
	} `toml:"auth"`
	Bot struct {
		Prefixes     []string `toml:"prefixes"`
		BotOwners    []string `toml:"bot_owners"`
		Invite       string   `toml:"invite"`
		CustomStatus struct {
			Override bool   `toml:"override"`
			Status   string `toml:"status"`
		} `toml:"custom_status"`
	} `toml:"bot"`
}

func getConfig() (config botConfig, err error) {
	configFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(configFile, &config)
	sugar.Infof("Loaded configuration file.")
	return config, err
}
