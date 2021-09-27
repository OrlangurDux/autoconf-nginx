package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

var (
	//SysConfig application config
	SysConfig Config
)

type Telegram struct {
	//ChatID caht id info from telegram
	ChatID int64
}

type Config struct {
	//Telegram telegram config section
	Telegram Telegram `yaml:"Telegram"`
}

func WriteSysConfig() {
	config, err := yaml.Marshal(&SysConfig)
	if err != nil {
		log.Print(err)
	}

	err = ioutil.WriteFile("config.yaml", config, 0666)
	if err != nil {
		log.Print(err)
	}

	log.Print("Config saved!")
}
