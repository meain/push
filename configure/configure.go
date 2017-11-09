package configure

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Conf has the settings configured by the user
type Conf struct {
	Token  string `yaml:"token"`
	Device string `yaml:"device"`
}

// GetConf parses the config file and get the user config
func (c *Conf) GetConf(file string) *Conf {
	file = expandFilePath("~/.push.yaml")
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
