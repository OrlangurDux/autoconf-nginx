package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	//variables scan directory for config yaml
	entryDir string
	//variables for save generated files conf
	vhostsDir string
)

// structure yaml file
type confNginx struct {
	Conf struct {
		Host      string `yaml:"host"`
		Container string `yaml:"container"`
		Port      int64
		Ssl       int64
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
		if err == nil && info.Name() == "nginx.yaml" {
			log.Print(info)
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
	flag.StringVar(&entryDir, "entry-dir", "", "Scan config directory")
	flag.StringVar(&vhostsDir, "vhosts-dir", "", "Scan config directory")
	flag.Parse()

	if entryDir == "" {
		log.Print("-entry-dir is required param")
		os.Exit(1)
	}
	if vhostsDir == "" {
		log.Print("-vhosts-dir is required param")
		os.Exit(1)
	}
}

func main() {
	configs := searchConfigFiles(entryDir)
	if configs == nil {
		log.Print("Configs files not found")
		os.Exit(1)
	} else {
		for _, f := range configs {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
					log.Print(err)
			}
			var template string = dir + "/template/"
			c, err := readConf(f)
			if err != nil {
				log.Fatal(err)
			}

			configFile := vhostsDir + c.Conf.Host + ".conf"
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				if c.Conf.Ssl == 1 {
					template = template + "nginx.ssl.template"
				} else {
					template = template + "nginx.nonssl.template"
				}

				file, err := ioutil.ReadFile(template)

				if err != nil {
					log.Print(err)
					os.Exit(1)
				}

				replace := bytes.Replace(file, []byte("#domain#"), []byte(c.Conf.Host), -1)
				replace = bytes.Replace(replace, []byte("#port#"), []byte(strconv.Itoa(int(c.Conf.Port))), -1)
				replace = bytes.Replace(replace, []byte("#container#"), []byte(c.Conf.Container), -1)

				if err := ioutil.WriteFile(configFile, replace, 0666); err != nil {
					log.Print(err)
					os.Exit(1)
				}

				fmt.Printf("%v", c)
			}
		}
	}
}
