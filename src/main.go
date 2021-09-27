package main

import (
	"autoconf/bot"
	"autoconf/config"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	//variables scan directory for config yaml
	entryDir string
	//variables for save generated files conf
	vhostsDir string
	//variable for configuration name file
	configName string
	//BotToken token for telegram bot
	botToken string
	//JoinKey key for join telegram chat
	joinKey string
)

// structure yaml file
type configHost struct {
	Conf struct {
		Host        string `yaml:"host"`
		Container   string `yaml:"container"`
		Port        int64
		Ssl         int64
		SslNameCert string `yaml:"sslNameCert"`
		SslNameKey  string `yaml:"sslNameKey"`
	}
}

//read yaml configuration to structure config host
func readHostConfig(filename string) (*configHost, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &configHost{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

//read yaml configuration to structure system configuration
func readSysConfig(filename string) (*config.Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &config.Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

//serach recrusive direcotry and subdirectory config file
func searchConfigFiles(path string) []string {
	var configs []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == configName {
			configs = append(configs, path)
		}
		return nil
	})
	if err != nil {
		if bot.Bot != nil {
			bot.SendBotMessage("Don't walk in dir " + path + ". Service is down. Error: " + err.Error())
		}
		log.Fatal(err)
	}
	return configs
}

// initializate param from cli
func init() {
	flag.StringVar(&entryDir, "input-dir", "", "Scan config directory")
	flag.StringVar(&vhostsDir, "output-dir", "", "Scan config directory")
	flag.StringVar(&configName, "config-name", "nginx.yaml", "Config name for search yaml")
	flag.StringVar(&botToken, "bot-token", "", "Telegram bot token")
	flag.StringVar(&joinKey, "join-key", "", "Key for join telegram chat")
	flag.Parse()

	if entryDir == "" {
		log.Print("-input-dir is required param")
		os.Exit(1)
	}
	if vhostsDir == "" {
		log.Print("-output-dir is required param")
		os.Exit(1)
	}

	if botToken != "" {
		if joinKey == "" {
			log.Print("-join-key is required param if use -bot-token param")
		} else {
			bot.BotToken = botToken
			bot.JoinKey = joinKey
		}
	}

	sysConfig, err := readSysConfig("config.yaml")
	if err != nil {
		config.SysConfig = config.Config{Telegram: config.Telegram{ChatID: 0}}
		log.Print(err)
	} else {
		config.SysConfig = *sysConfig
		bot.ChatID = sysConfig.Telegram.ChatID
	}
}

func generateFileConfig() {
	for {
		configs := searchConfigFiles(entryDir)
		if configs == nil {
			log.Print("Configs files not found. Sleep...")
		} else {
			bServiceRestart := false
			for _, f := range configs {
				dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					if bot.Bot != nil {
						bot.SendBotMessage("Don't read dir " + dir + ". Error: " + err.Error())
					}
					log.Print(err)
				}
				var template string = dir + "/template/"
				c, err := readHostConfig(f)
				if err != nil {
					if bot.Bot != nil {
						bot.SendBotMessage("Don't parse config file " + f + ". Error:" + err.Error())
					}
					log.Print(err)
				}

				bUpdate := false

				configFile := vhostsDir + c.Conf.Host + ".conf"
				_, err = os.Stat(configFile)
				if err == nil {
					modifiedConfigFile, err := os.Stat(configFile)
					if err != nil {
						if bot.Bot != nil {
							bot.SendBotMessage("Don't get modified date from config file. Error: " + err.Error())
						}
						log.Print(err)
					}
					modifiedConfigTime := modifiedConfigFile.ModTime().Unix()

					modifiedYamlFile, err := os.Stat(f)

					if err != nil {
						if bot.Bot != nil {
							bot.SendBotMessage("Don't get modified date from yaml file. Error: " + err.Error())
						}
						log.Print(err)
					}

					modifiedYamlTime := modifiedYamlFile.ModTime().Unix()

					if modifiedYamlTime >= modifiedConfigTime {
						bUpdate = true
					}
				}
				if os.IsNotExist(err) || bUpdate {
					if c.Conf.Ssl == 1 {
						template = template + "nginx.ssl.template"
					} else {
						template = template + "nginx.nonssl.template"
					}

					file, err := ioutil.ReadFile(template)

					if err != nil {
						if bot.Bot != nil {
							bot.SendBotMessage("Don't read template file. Error: " + err.Error())
						}
						log.Print(err)
					} else {

						replace := bytes.Replace(file, []byte("#domain#"), []byte(c.Conf.Host), -1)
						replace = bytes.Replace(replace, []byte("#port#"), []byte(strconv.Itoa(int(c.Conf.Port))), -1)
						replace = bytes.Replace(replace, []byte("#container#"), []byte(c.Conf.Container), -1)
						replace = bytes.Replace(replace, []byte("#sslnamecert#"), []byte(c.Conf.SslNameCert), -1)
						replace = bytes.Replace(replace, []byte("#sslnamekey#"), []byte(c.Conf.SslNameKey), -1)

						if err := ioutil.WriteFile(configFile, replace, 0666); err != nil {
							if bot.Bot != nil {
								bot.SendBotMessage("Don't write config file " + configFile + ". Error: " + err.Error())
							}
							log.Print(err)
						} else {
							var message string
							if bUpdate {
								message = "Configuration file " + configFile + " updated."
							} else {
								message = "Configuration file " + configFile + " added."
							}
							if bot.Bot != nil {
								bot.SendBotMessage(message)
							}
							bServiceRestart = true
						}
					}
				}

				if bServiceRestart {
					cmd := exec.Command("/bin/sh", "-c", "/etc/init.d/nginx check-reload")

					stdout, err := cmd.CombinedOutput()

					if err != nil {
						if bot.Bot != nil {
							bot.SendBotMessage(string(stdout[:]))
						}
						log.Print(string(stdout[:]))
						fmt.Println(err)
						//os.Exit(1)
					}
				}
			}

			time.Sleep(time.Millisecond * time.Duration(10000))
		}
	}
}

func main() {
	go generateFileConfig()
	bot.InitBot()
}
