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
	ChatID int64 `yaml:"chat_id"`
	//JoinKey key for hoin bot channel
	JoinKey string `yaml:"join_key"`
	//BotToken telegram bot token
	BotToken string `yaml:"bot_token"`
}

type Settings struct {
	//EntryDir input dir
	EntryDir string `yaml:"input_dir"`
	//VhostsDir output dir
	VhostsDir string `yaml:"output_dir"`
	//ConfigName name input config file default "nginx.yaml"
	ConfigName string `yaml:"config_name"`
}

type Config struct {
	//Telegram telegram config section
	Telegram Telegram `yaml:"telegram"`
	//Settings config section
	Settings Settings `yaml:"settings"`
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
