package main

import (
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

	"autoconf/bot"

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
type confNginx struct {
	Conf struct {
		Host        string `yaml:"host"`
		Container   string `yaml:"container"`
		Port        int64
		Ssl         int64
		SslNameCert string `yaml:"sslNameCert"`
		SslNameKey  string `yaml:"sslNameKey"`
	}
}

//read yaml configuration to structure coonfNginx
func readConf(filename string) (*confNginx, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &confNginx{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

//serach recrusive direcotry and subdirectory config file
func searchConfigFiles(path string) []string {
	var configs []string
	var e = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == configName {
			configs = append(configs, path)
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
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
					log.Print(err)
				}
				var template string = dir + "/template/"
				c, err := readConf(f)
				if err != nil {
					log.Print(err)
				}

				bUpdate := false

				configFile := vhostsDir + c.Conf.Host + ".conf"
				_, err = os.Stat(configFile)
				if err == nil {
					modifiedConfigFile, err := os.Stat(configFile)
					if err != nil {
						log.Print(err)
					}
					modifiedConfigTime := modifiedConfigFile.ModTime().Unix()

					modifiedYamlFile, err := os.Stat(f)

					if err != nil {
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
						log.Print(err)
					} else {

						replace := bytes.Replace(file, []byte("#domain#"), []byte(c.Conf.Host), -1)
						replace = bytes.Replace(replace, []byte("#port#"), []byte(strconv.Itoa(int(c.Conf.Port))), -1)
						replace = bytes.Replace(replace, []byte("#container#"), []byte(c.Conf.Container), -1)
						replace = bytes.Replace(replace, []byte("#sslnamecert#"), []byte(c.Conf.SslNameCert), -1)
						replace = bytes.Replace(replace, []byte("#sslnamekey#"), []byte(c.Conf.SslNameKey), -1)

						if err := ioutil.WriteFile(configFile, replace, 0666); err != nil {
							log.Print(err)
						} else {
							bServiceRestart = true
						}
					}
				}

				if bServiceRestart {
					cmd := exec.Command("/bin/sh", "-c", "/etc/init.d/nginx check-reload")

					stdout, err := cmd.CombinedOutput()

					if err != nil {
						//bot.SendBotMessage("Test")
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
